package mistral_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mistral "github.com/psyb0t/go-mistral"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAPIKey      = "test-key"
	testModel       = "mistral-small-latest"
	testReply       = "hello back"
	modelsPath      = "/v1/models"
	completionsPath = "/v1/chat/completions"
	headerAuth      = "Authorization"
	headerType      = "Content-Type"
	contentTypeJSON = "application/json"
)

// recordedRequest is what the fake server saw.
type recordedRequest struct {
	method string
	path   string
	auth   string
	body   []byte
}

func newServer(
	t *testing.T,
	responseBody string,
	seen *recordedRequest,
) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			seen.method = r.Method
			seen.path = r.URL.Path
			seen.auth = r.Header.Get(headerAuth)

			var body json.RawMessage
			if r.ContentLength > 0 {
				assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			}

			seen.body = body

			w.Header().Set(headerType, contentTypeJSON)
			_, err := w.Write([]byte(responseBody))
			assert.NoError(t, err)
		},
	))
	t.Cleanup(server.Close)

	return server
}

func newClient(
	t *testing.T,
	server *httptest.Server,
) *mistral.ClientWithResponses {
	t.Helper()

	auth := mistral.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set(headerAuth, "Bearer "+testAPIKey)

			return nil
		},
	)

	client, err := mistral.NewClientWithResponses(server.URL, auth)
	require.NoError(t, err)

	return client
}

func TestListModels(t *testing.T) {
	t.Parallel()

	var seen recordedRequest

	server := newServer(t, `{
		"object": "list",
		"data": [{"id": "`+testModel+`", "object": "model"}]
	}`, &seen)

	res, err := newClient(t, server).ListModelsV1ModelsGetWithResponse(
		t.Context(),
		&mistral.ListModelsV1ModelsGetParams{},
	)
	require.NoError(t, err)

	assert.Equal(t, http.MethodGet, seen.method)
	assert.Equal(t, modelsPath, seen.path)
	assert.Equal(t, "Bearer "+testAPIKey, seen.auth)

	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	require.NotNil(t, res.JSON200.Data)
	require.Len(t, *res.JSON200.Data, 1)

	card, err := (*res.JSON200.Data)[0].AsBaseModelCard()
	require.NoError(t, err)
	assert.Equal(t, testModel, card.Id)
}

func TestChatCompletion(t *testing.T) {
	t.Parallel()

	var seen recordedRequest

	server := newServer(t, `{
		"id": "cmpl-1",
		"object": "chat.completion",
		"created": 1,
		"model": "`+testModel+`",
		"usage": {
			"prompt_tokens": 3,
			"completion_tokens": 2,
			"total_tokens": 5
		},
		"choices": [{
			"index": 0,
			"finish_reason": "stop",
			"message": {"role": "assistant", "content": "`+testReply+`"}
		}]
	}`, &seen)

	var content mistral.UserMessage_Content
	require.NoError(t, content.FromUserMessageContent0("hello"))

	role := mistral.UserMessageRoleUser

	var message mistral.ChatCompletionRequest_Messages_Item
	require.NoError(t, message.FromUserMessage(mistral.UserMessage{
		Content: &content,
		Role:    &role,
	}))

	res, err := newClient(t, server).
		ChatCompletionV1ChatCompletionsPostWithResponse(
			t.Context(),
			mistral.ChatCompletionRequest{
				Model: testModel,
				Messages: []mistral.ChatCompletionRequest_Messages_Item{
					message,
				},
			},
		)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, seen.method)
	assert.Equal(t, completionsPath, seen.path)
	assert.Equal(t, "Bearer "+testAPIKey, seen.auth)
	assert.JSONEq(t, `{
		"model": "`+testModel+`",
		"messages": [{"role": "user", "content": "hello"}]
	}`, string(seen.body))

	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	require.Len(t, res.JSON200.Choices, 1)

	choice := res.JSON200.Choices[0]
	assert.Equal(t, mistral.ChatCompletionChoiceFinishReasonStop,
		choice.FinishReason)
	require.NotNil(t, choice.Message.Content)

	reply, err := choice.Message.Content.AsAssistantMessageContent0()
	require.NoError(t, err)
	assert.Equal(t, testReply, reply)
}
