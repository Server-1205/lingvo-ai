[← Конфигурация](configuration.md) · [Back to README](../README.md)

# Streaming AI Responses

The chat endpoint supports Server-Sent Events (SSE) for real-time streaming of AI responses.

## Endpoint

`POST /api/chat/stream`

Same authentication and rate-limiting as `POST /api/chat`.

## Request

```json
{"text": "Hello, how are you?"}
```

## Response Format

The response uses `Content-Type: text/event-stream`. Each line is prefixed with `data: ` followed by a JSON event:

```
data: {"type":"token","data":"Hello"}
data: {"type":"token","data":"! I"}
data: {"type":"token","data":"'m doing great!"}
data: {"type":"corrections","data":[{"original":"im","corrected":"I'm","explanation_uz":"...","explanation_ru":"...","type":"grammar"}]}
data: {"type":"usage","data":{"daily_used":1,"daily_limit":10,"is_premium":false}}
data: {"type":"done"}
```

### Event Types

| Type | Payload | Description |
|------|---------|-------------|
| `token` | `string` | Partial text chunk from Gemini |
| `corrections` | `Correction[]` | Grammar/vocabulary corrections (sent once) |
| `usage` | `Usage` | Updated daily usage info |
| `done` | *(none)* | Signals stream completion |

## Backend Architecture

1. `internal/ai/gemini.go` — `ChatStream()` method uses `GenerateContentStream()` from the Gemini SDK, emits tokens via a channel
2. `internal/api/chat_stream.go` — SSE handler reads the channel and writes events to the response writer
3. `internal/ai/response.go` — `StreamEvent` type definition with `token` and `result` variants

## Frontend Integration

The React Chat component uses `fetch` with `ReadableStream` to consume the SSE stream:

- `web/src/api/client.ts` — `chatStream()` function parses SSE lines and dispatches events via callbacks
- `web/src/components/Chat.tsx` — replaces the previous `useMutation` approach with real-time token rendering

## See Also

- [API Reference](api.md) — все эндпоинты
- [Установка и запуск](getting-started.md) — как запустить проект
