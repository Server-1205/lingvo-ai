# Implementation Plan: Streaming AI Responses

Branch: feature/streaming-ai-responses
Created: 2026-06-06

## Settings
- Testing: yes
- Logging: verbose
- Docs: yes

## Roadmap Linkage
Milestone: "Асинхронный streaming ответов" (Этап 1: AI + Чат)
Rationale: Enables real-time AI response streaming via SSE, replacing the current blocking request-response pattern

## Commit Plan
- **Commit 1** (after tasks 1-2): "feat(ai): add streaming client method and stream event types"
- **Commit 2** (after tasks 3-4): "feat(api): add SSE streaming endpoint for /api/chat/stream"
- **Commit 3** (after tasks 5-6): "feat(web): consume streaming API with real-time token display"

## Tasks

### Phase 1: Backend — Streaming AI Client

- [x] Task 1: Add `ChatStream` method and stream event types to AI client

  Add to `internal/ai/gemini.go`:
  - `StreamEvent` struct with `Type` (string: "token"/"corrections"/"usage") and `Data` (json.RawMessage)
  - `ChatStream(ctx, prompt) <-chan StreamEvent` method using `model.GenerateContentStream()`
  - Each `GenerateContentResponseChunk` yields a `{type:"token", data:"..."}` event
  - After stream completes, parse the full text as JSON → extract corrections → emit `{type:"corrections"}` and `{type:"usage"}` events, then close channel

  LOGGING:
  - DEBUG: each chunk received from Gemini (`chunk received, %d bytes`)
  - INFO: stream started / completed
  - ERROR: if Gemini returns an error mid-stream

  Files: `internal/ai/gemini.go`, `internal/ai/response.go` (if StreamEvent types go there)

- [x] Task 2: Write tests for `ChatStream`

  Create `internal/ai/gemini_test.go`:
  - Test that stream emits token events (mock a few chunks)
  - Test that stream emits corrections event after full response is accumulated
  - Test that stream closes on context cancellation

  Files: `internal/ai/gemini_test.go`

### Phase 2: Backend — SSE Endpoint
<!-- Commit checkpoint: tasks 1-2 -->

- [x] Task 3: Create `POST /api/chat/stream` SSE handler

  Create `internal/api/chat_stream.go`:
  - Accept same `models.ChatRequest` body as `/api/chat`
  - Read user_id, lang, level, is_premium from context (same middleware chain)
  - Build prompt via `ai.BuildChatPrompt()`
  - Call `aiClient.ChatStream()` with request context
  - Set response headers: `Content-Type: text/event-stream`, `Cache-Control: no-cache`, `Connection: keep-alive`
  - Flush after each write
  - On each stream event: write `data: <json>\n\n` to response writer
  - Increment message count in DB after stream completes (using `db.IncrementMessageCount`)
  - Handle client disconnect (`c.Request.Context().Done()`) to cancel the Gemini stream

  LOGGING:
  - INFO: streaming session started / completed for user
  - DEBUG: each event written
  - ERROR: Gemini errors, DB errors, serialize errors
  - WARN: early client disconnect

  Files: `internal/api/chat_stream.go`

- [x] Task 4: Register streaming route and update router

  Add to `internal/api/router.go`:
  - `protected.POST("/chat/stream", rateMw, chatStreamHandler)` — note: rate limit applied (streaming counts as one message regardless of length)

  Files: `internal/api/router.go`
  <!-- Commit checkpoint: tasks 3-4 -->

### Phase 3: Frontend — Streaming

- [x] Task 5: Add streaming API client function

  Add to `web/src/api/client.ts`:
  - `chatStream(req: ChatRequest, onToken: (text: string) => void, onCorrections: (corrections: Correction[]) => void, onUsage: (usage: Usage) => void, onDone: () => void): Promise<void>`
  - Uses `fetch` with POST to `/api/chat/stream`
  - Reads response body via `response.body.getReader()` (ReadableStream)
  - Decodes `TextDecoder` and parses lines starting with `data: `
  - For each JSON event: call appropriate callback
  - Type definitions: `StreamEvent` union type

  LOGGING:
  - DEBUG: each event received

  Files: `web/src/api/client.ts`

- [x] Task 6: Update Chat component for streaming

  Update `web/src/components/Chat.tsx`:
  - Replace `useMutation` with a streaming approach
  - On send: create a placeholder AI message with empty text, then update it via `onToken` callback as tokens arrive
  - Show corrections once `onCorrections` fires
  - Show usage via `onUsage`
  - Handle loading state during streaming
  - Handle errors gracefully

  Note: Keep the non-streaming fallback for backward compat, or replace entirely — prefer replacing entirely since SSE endpoint has the same middleware/semantics.

  Files: `web/src/components/Chat.tsx`
  <!-- Commit checkpoint: tasks 5-6 -->
