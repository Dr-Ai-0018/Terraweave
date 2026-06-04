# TerraWeave API

Go backend for TerraWeave.

## Current Scope

- OpenAI-compatible `POST /v1/responses` proxy.
- Streaming-only transport. The backend always forwards `stream: true`.
- Custom provider base URL through `TERRAWEAVE_AI_BASE_URL`.
- Custom provider API key through `TERRAWEAVE_AI_API_KEY`.
- Custom model list through `TERRAWEAVE_AI_MODELS`.
- Default model: `gpt-5.4`.
- Modal GPU option registry with default GPU `L4`.
- Modal GPU prices are copied from Modal's public pricing page and should be refreshed before serious cost-sensitive runs.

The API does not use an OpenAI SDK. It uses direct HTTP requests so OpenAI-compatible providers can be configured by URL and key.

## Endpoints

```text
GET  /healthz
GET  /v1/models
GET  /v1/modal/gpus
POST /v1/responses
```

## Streaming Responses

`POST /v1/responses` expects an OpenAI Responses-style JSON body and returns `text/event-stream`.

The request body is passed through to the configured provider with these changes:

- `stream` is forced to `true`.
- `model` defaults to `TERRAWEAVE_AI_DEFAULT_MODEL` when omitted.

Tool/function payloads are passed through. TerraWeave will add local tool execution, MCP tool discovery, and skill-backed tool registration above this streaming loop instead of depending on a vendor SDK.

Non-streaming calls are intentionally not supported. If a caller sends `stream: false` or omits `stream`, the backend overrides it to `true`.

## Environment

```text
TERRAWEAVE_API_ADDR=:8080
TERRAWEAVE_AI_BASE_URL=https://api.openai.com
TERRAWEAVE_AI_API_KEY=
TERRAWEAVE_AI_DEFAULT_MODEL=gpt-5.4
TERRAWEAVE_AI_MODELS=gpt-5.4
TERRAWEAVE_MODAL_DEFAULT_GPU=L4
```

## Run

```bash
go run ./cmd/server
```

When running inside the Codex sandbox, set a writable build cache:

```bash
GOCACHE=/tmp/terraweave-go-cache go test ./...
```

## Modal GPU Options

Default: `L4`.

Known public options:

```text
T4
L4
A10
A10G
L40S
A100
A100-40GB
A100-80GB
RTX-PRO-6000
H100
H100!
H200
B200
B200+
```

`B200+` is Modal's documented opt-in path for B200 or B300-compatible capacity. `A10G` is included because Modal's playground and A10G article use it directly, even though the main GPU guide lists `A10`.
