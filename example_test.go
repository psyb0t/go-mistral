package mistral_test

import (
	"context"
	"net/http"
	"os"

	mistral "github.com/psyb0t/go-mistral"
)

// The README usage example. It calls the live API, so it is compiled but not
// run.
//
//nolint:testableexamples // needs a real API key and network access
func Example() {
	auth := mistral.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set(
				"Authorization",
				"Bearer "+os.Getenv("MISTRAL_API_KEY"),
			)

			return nil
		},
	)

	client, err := mistral.NewClientWithResponses(
		"https://api.mistral.ai",
		auth,
	)
	if err != nil {
		panic(err)
	}

	res, err := client.ListModelsV1ModelsGetWithResponse(
		context.Background(),
		&mistral.ListModelsV1ModelsGetParams{},
	)
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
