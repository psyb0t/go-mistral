# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.5.0, 2026-11-26

- Breaking: the RAG deployment operations are no longer generated. `GetDeploymentSummariesV1RagDeploymentsGet`, `RegisterDeploymentV1RagDeploymentsPut`, `UnregisterDeploymentV1RagDeploymentsDeploymentIdDelete` and `UpdateIndexMetricsV1RagDeploymentsDeploymentIdMetricsPut` are removed from `Client` and `ClientWithResponses`, with their request, response and `RegisterDeploymentRequest*`, `GetDeploymentSummariesResponse*` and `UpdateMetricsRequest*` types.
- Breaking: `PostDatasetRecordsFromPlaygroundV1ObservabilityDatasetsDatasetIdImportsFromPlaygroundPost` and `ImportDatasetFromPlaygroundRequest` are removed. `CreateIngestionPipelineConfigurationRequest.TargetIndexes`, `IngestionPipelineTargetIndexRef` and the `SchemaField*` types are removed.
- Breaking: `ModerationsV1ModerationsPostJSONRequestBody` changed from `ClassificationRequest` to `ModerationRequest`. `UsersApiAdminWorkspacesAddOrUpdateUsersWorkspacesJSONRequestBody` changed from `AdminWorkspaceMemberIN` to `AdminWorkspaceMemberUpsertIN`.
- Breaking: `FromToolChoice`, `MergeToolChoice` and `AsToolChoice` on `ChatCompletionRequest_ToolChoice` and `AgentsCompletionRequest_ToolChoice` are replaced by `FromNamedToolChoice` and `MergeNamedToolChoice`. The `ToolChoice` type is removed.
- Breaking: `EmbeddingRequestInput1` and `EmbeddingRequest_Input.FromEmbeddingRequestInput1` are removed. `FromAudioChunk`, `FromImageURLChunk`, `FromEmbeddingRequestInput3` and `FromEmbeddingRequestInput4` are added, with matching `Merge` methods.
- Breaking: several role enum constants are removed, such as `AdminWorkspaceMemberSingleINRoleNameBilling` and `AdminAssignGroupToWorkspaceInRoleNamesMember`, along with the `StaticWorkspaceRoles*` hash-named types. `SummaryMetadata.Status` changed to `SummaryStatus1`, and `EvaluationRunV2Response.RunScores` and `EvaluationRunV2WithRefsResponse.RunScores` now use `EvaluationRunRunScoreV2`.
- Breaking: `SpeechStreamEvents`, `SpeechStreamEvents_Data` and `SummaryStreamTypes` are removed.
- Added `CreateDatasetRecordsV1ObservabilityDatasetsDatasetIdRecordsBulkPost`, `CreateDeploymentArchiveV1WorkflowsDeploymentsArchivesPost` and `UsersApiAdminWorkspacesListWorkspaceUsers`, each with `WithBody` and `WithResponse` variants where applicable.
- Added `AddRoleNames`, `AddRoles`, `RemoveRoleNames` and `RemoveRoles` on `AdminUpdateGroupWorkspaceAssignmentIn`, and turn detection event unions on `TranscriptionStreamEvents_Data` such as `FromTurnDetectionStartEvent` and `FromTurnDetectionEndEvent`.

## v0.4.0, 2026-11-18

- Breaking: the Judges API is no longer generated. `CreateJudgeV1ObservabilityJudgesPost`, `GetJudgesV1ObservabilityJudgesGet`, `GetJudgeByIdV1ObservabilityJudgesJudgeIdGet`, `UpdateJudgeV1ObservabilityJudgesJudgeIdPut`, `DeleteJudgeV1ObservabilityJudgesJudgeIdDelete` and the two live judging operations are removed from `Client` and `ClientWithResponses`, with the `Judge`, `JudgeOutput`, `CreateJudgeRequest`, `JudgeConversationRequest`, `JudgeDatasetRecordRequest` and `ListJudgesResponse` types.
- Breaking: `UsersApiAdminUserGroupsGetNestedGroups` and `UsersApiAdminUserGroupsSetNestedGroups` are removed, with `NestedGroupRef`, `NestedGroupsOut` and the `ParentGroupIds` fields on `AdminUserGroupIn`, `AdminUpdateUserGroupIn` and `AdminUserGroupOut`.
- Breaking: `RawRoles` on `AdminOrganizationMemberOUT`, `AdminUserOUT`, `MemberWorkspaceInfo` and `OrganizationInviteOUT`, and `Roles` on `OrganizationUserInviteOUT`, are now plain slices such as `[]StaticOrganizationRoles` instead of union wrapper types. The `*_RawRole` and `*_RawRoles` wrapper types are removed.
- Breaking: the `Filler` field on `HangUpTool`, `TransferCallTool` and the `VoiceAgent*Tool` types changed from `*StaticToolFiller` to an inline union type. `ConnectorToolResultMetadata.StructuredContent` changed from `*map[string]interface{}` to `interface{}`. `CreateWebhook.EventTypePath` and the constant `ConnectorCreateCredentialsV1ParamsConsumerScopeOrganization` are removed.
- Added scorer operations: `ListScorersV1ObservabilityScorersGet`, `CreateScorerV1ObservabilityScorersPost`, `GetScorerV1ObservabilityScorersScorerSlugGet` and `EditScorerV1ObservabilityScorersScorerSlugPatch`.
- Added `UpdateSessionV1AudioSessionsSessionIdPatch`, `UsersApiAdminUsersGetUserWorkspaces` and `UsersApiAdminUsersReactivateUser`, each with a `WithResponse` variant.
- Added `LastLoginAt`, `Status` and `WorkspacesTotal` on `AdminUserOUT` and `AdminOrganizationMemberOUT`, `BridgeMetadata` on `AudioSessionResponse` and `AudioAgentSessionSummary`, the `AmazonConnectContact` type, and `AggregationRowMetricValue2` for `AggregationRow_MetricValue`.
- Raised the Go requirement to 1.26.11 in `go.mod`, `Dockerfile.dev` and the setup guide, to pick up the fix for GO-2026-6859 in `net/http`.

## v0.3.0, 2026-11-05

- Development image package installation uses HTTPS and bounded retries for signed snapshot indexes and package downloads.
- Breaking: regenerated Mistral specification 1.0.0. Legacy managed-index operations and several earlier response models are no longer generated. Consumers of those names must update against the current generated API.

## v0.2.0, 2026-10-29

- The client is generated from the spec Mistral AI publishes at `https://docs.mistral.ai/openapi.yaml`. The platform-docs-public repository it was downloaded from before no longer exists, so scheduled regeneration had stopped.
- Breaking: the current spec has 317 operations instead of 169. It adds 176, including the Responses API, workflow deployments, managed RAG indexes, service accounts, connector sharing and observability traces and spans. It removes 28: the fine-tuning job routes, the observability campaign and chat completion event routes, and the connector credential routes per organization, workspace and user, which are replaced by `ConnectorCreateCredentialsV1`, `ConnectorUpdateCredentials` and `ConnectorDeleteCredentials`. Many model types are renamed or removed with them.
- `internal/cmd/patch-spec` moves the schemas the spec nests in local `$defs` blocks into `components/schemas`, since kin-openapi cannot resolve `#/$defs/...` references. When two query parameters of one operation would get the same Go field name, the later one gets a numeric suffix, so the `sort_direction` alias of `sort.direction` is `SortDirection2`.

## v0.1.5, 2026-10-28

- The Semgrep scan leaves out `client.gen.go`, as its comment already said it should. `--exclude` did not apply because the file was named on the command line, and Semgrep timed out on it, so its findings came and went with the runner's speed and failed the scan at random.
- Updated Semgrep in the development image to 1.182.1 and refreshed its locked dependencies. certifi 2026.10.23 drops the emSign root CAs whose intermediate keys were compromised (GHSA-6q48-6c7w-9gw2).
- Dependabot no longer opens pull requests for `mcp` or the OpenTelemetry packages in `requirements-semgrep.txt`. Semgrep pins their versions itself, so bumping one of them on its own cannot install.
- Updated kin-openapi, which the code generator uses, from 0.142.0 to 0.149.0. The newer version fixes several panics on malformed specs reported in GitHub security advisories.
- The security scan merges its two gitleaks runs into one SARIF run. GitHub code scanning rejects a report with two runs from the same tool, so no scan results had been uploaded.

## v0.1.4, 2026-10-27

- The development image is based on Debian (`golang:1.26.10-trixie`) instead of Alpine. Alpine keeps only the latest build of each package, so the exact pins broke the image build whenever a package was updated. System packages now come from the snapshot.debian.org date the base image was built from, which keeps them fixed for a given base image digest. jq 1.8.2, shellcheck 0.11.0 and shfmt 3.13.1 keep their versions, and Semgrep runs in its own virtual environment.
- Raised the minimum Go version in `go.mod` from 1.26.9 to 1.26.10, which fixes the `crypto/x509` vulnerability GO-2026-6809 found by govulncheck.
- The issue-pull workflow runs at hours derived from the repository name instead of the same four hours as every other repository.

## v0.1.3, 2026-10-17

- Documentation only. The agent skill explains how to look up every endpoint, parameter and model with `go doc` or on pkg.go.dev, and declares its license.
- The skill no longer lists the API credentials under OpenClaw's `requires.env`, which kept OpenClaw from loading it until they were set. It requires only `go`; the credentials are still documented as what your program reads.
- The README points to the full API reference, and the examples in the README and the skill were checked against the current code.

## v0.1.2, 2026-10-17

- Documentation only. The README and `AGENTS.md` describe the daily regeneration workflow: what it checks, how it fixes failing checks, how it picks the next version and which paths a fix may not touch.

## v0.1.1, 2026-10-16

- Added the `go-mistral` agent skill. It installs with `claude plugin install go-mistral@psyb0t` in Claude Code, `codex plugin add go-mistral@psyb0t` in Codex and `openclaw skills install @psyb0t/go-mistral` from ClawHub.

## v0.1.0, 2026-10-16

- Added `mistral`, a REST client generated with oapi-codegen from Mistral AI's OpenAPI document (spec version 1.0.0).
- `make generate` downloads the current spec from Mistral AI's platform-docs-public repository and regenerates the client. Boolean `const` values are dropped before generation, since one of them collides with a component schema name. The spec is not committed.
