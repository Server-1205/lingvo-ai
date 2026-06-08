# TTS произношение (озвучка слов)

**Branch:** `feature/tts-pronunciation`
**Created:** 2026-06-07 11:15

## Settings

- **Testing:** yes — unit tests for TTS client + handler
- **Logging:** verbose — DEBUG logs for TTS requests, timing, cache hits
- **Docs:** yes — API endpoint doc + configuration guide

## Roadmap Linkage

- **Milestone:** "Фичи V2"
- **Rationale:** TTS произношение — запланированная V2 фича (ROADMAP.md:154)

## Research Context

- TTS через edge-tts (Python, бесплатно, Microsoft Edge TTS движок)
- Встроенные голоса: uz-UZ-MadinaNeural (узбекский), ru-RU-SvetlanaNeural (русский)
- Вызов через subprocess: `edge-tts --voice <voice> --text <text> --write-media <file>`
- На фронтенде: нативный `Audio()` API, без новых зависимостей

---

## Tasks

### Phase 1: Backend TTS

#### Task 1: [x] Создать `internal/tts/` пакет с TTS-клиентом

**Файлы:** `internal/tts/client.go`, `internal/tts/client_test.go`

**Что сделать:**
- Создать пакет `internal/tts/`
- `Client` struct с полями: `voiceUz`, `voiceRu`, `pythonPath` (путь к `edge-tts`), `sugar *zap.SugaredLogger`
- `NewClient(voiceUz, voiceRu string, sugar *zap.SugaredLogger) *Client`
- `func (c *Client) Synthesize(ctx context.Context, text, lang string) ([]byte, error)`
  - Выбрать голос по lang: `uz` → `c.voiceUz`, `ru` → `c.voiceRu`
  - Создать временный файл через `os.CreateTemp`
  - Выполнить `exec.CommandContext(ctx, "edge-tts", "--voice", voice, "--text", text, "--write-media", tmpFile)`
  - Прочитать файл, удалить, вернуть `[]byte`
  - Логи: `[tts] starting synthesis`, `[tts] synthesized X bytes in Yms`
  - Если edge-tts не установлен — вернуть понятную ошибку
- `func (c *Client) IsAvailable() bool` — проверка наличия edge-tts

**Тесты:**
- `TestSynthesize_Success` — mock exec.Command, проверить что возвращаются байты
- `TestSynthesize_Fail` — команда падает с ошибкой
- `TestSynthesize_EmptyText` — пустой текст
- `TestVoiceSelection` — правильный голос для uz/ru

#### Task 2: [x] Создать `internal/api/tts.go` — хендлер TTS

**Файлы:** `internal/api/tts.go`, `internal/api/tts_test.go`

**Что сделать:**
- `ttsHandler(ttsClient *tts.Client, sugar *zap.SugaredLogger) gin.HandlerFunc`
- `GET /api/tts?text=...&lang=uz`
- Параметры:
  - `text` (required, max 500 символов)
  - `lang` (optional, default "uz")
- Валидация: text обязателен, не пустой, не длиннее 500
- Вызов `ttsClient.Synthesize(ctx, text, lang)`
- Ответ: `Content-Type: audio/mpeg`, `Cache-Control: no-store`, тело — аудио-байты
- При ошибке: `500 {"error": "tts_failed"}`
- Rate limit: проверять `c.Get("daily_used")` — TTS тоже расходует дневной лимит
- Логи: `[tts] request text=..., lang=..., size=...`

**Тесты:**
- `TestTTSHandler_Success` — мок TTS клиента, проверить Content-Type и тело
- `TestTTSHandler_MissingText` — 400
- `TestTTSHandler_TextTooLong` — 400 (500+ символов)
- `TestTTSHandler_Failure` — мок возвращает ошибку → 500

#### Task 3: Подключить TTS в `router.go` и `main.go`

**Файлы:** `internal/api/router.go`, `cmd/server/main.go`

**Что сделать:**
- `router.go`: добавить `tts.POST("/tts", ttsHandler(ttsClient, sugar))` — защищён auth middleware
- `router.go`: обновить сигнатуру `RegisterRoutes` — добавить параметр `ttsClient *tts.Client`
- `main.go`: инициализировать `tts.NewClient("uz-UZ-MadinaNeural", "ru-RU-SvetlanaNeural", sugar)`
- `main.go`: если TTS недоступен (`!ttsClient.IsAvailable()`), залогировать WARN, но не падать
- `main.go`: передать `ttsClient` в `RegisterRoutes`

#### Task 4: Добавить конфигурацию TTS

**Файлы:** `.env`, `docs/configuration.md`, `docs/api.md`

**Что сделать:**
- `.env`: добавить `TTS_VOICE_UZ=uz-UZ-MadinaNeural`, `TTS_VOICE_RU=ru-RU-SvetlanaNeural` (опционально, с дефолтами)
- `docs/configuration.md`: добавить секцию TTS с переменными
- `docs/api.md`: добавить `GET /api/tts` с параметрами, ответом, ошибками

---

### Phase 2: Frontend TTS

#### Task 5: Создать хук `useTTS`

**Файлы:** `web/src/hooks/useTTS.ts`

**Что сделать:**
- `function useTTS()` — возвращает `{ play: (text: string, lang: string) => Promise<void>, isPlaying: boolean, error: string | null }`
- Вызов `fetch('/api/tts?text=' + encodeURIComponent(text) + '&lang=' + lang)` с initData-заголовком
- Воспроизведение через `new Audio(url)` → `audio.play()`
- Очистка `URL.createObjectURL` после окончания воспроизведения
- Логи: `[tts] playing text=..., lang=...`
- Защита: блокировка повторного нажатия пока играет

#### Task 6: Добавить TTS-кнопку в Vocabulary.tsx

**Файлы:** `web/src/components/Vocabulary.tsx`

**Что сделать:**
- Добавить 🔊 кнопку в каждый word card (рядом с level badge, строка 158-162)
- Импортировать `useTTS` из `../hooks/useTTS`
- При нажатии: `play(w.word, 'uz')` — язык по умолчанию
- Пока играет: кнопка disabled, иконка 🔄
- Паттерн: `{!isLoading && <button ...>` (из skill-context rules)

#### Task 7: Добавить TTS-кнопку в ReviewCard.tsx

**Файлы:** `web/src/components/ReviewCard.tsx`

**Что сделать:**
- Добавить 🔊 кнопку на лицевой стороне карточки (рядом со словом) и на обратной стороне
- Импортировать `useTTS`
- При нажатии: `play(word.word, 'uz')`
- Кнопка с иконкой, без текста, компактная

#### Task 8: Добавить i18n ключи для TTS

**Файлы:** `web/src/locales/uz.json`, `web/src/locales/ru.json`

**Что сделать:**
- Добавить ключ `tts.play` = "O'qish" / "Озвучить"
- Добавить ключ `tts.playing` = "O'qilmoqda..." / "Воспроизводится..."

---

## Risks & Considerations

- **edge-tts не установлен** — сервер стартует с WARN, TTS эндпоинт возвращает 500. Документировать в конфигурации
- **Долгий синтез** — длинные тексты (>200 символов) могут синтезироваться несколько секунд. timeout в http.Client на фронтенде
- **Temp-файлы** — гарантировать удаление через `defer os.Remove`. Если сервер упал — `/tmp` очистится сам
- **Rate limit** — TTS использует тот же дневной лимит (10/day free, premium unlimited). Можно сделать отдельный лимит позже
- **Отсутствие голоса** — Microsoft может изменить голоса. При ошибке синтеза — fallback на генерацию через AI (Gemini TTS пока нет)

## Commit Plan

1. `feat(tts): add TTS client package (edge-tts)` — Task 1
2. `feat(tts): add TTS API endpoint` — Task 2
3. `feat(tts): wire TTS into router and server` — Task 3
4. `docs: add TTS configuration and API docs` — Task 4
5. `feat(tts): add useTTS hook` — Task 5
6. `feat(tts): add TTS button to vocabulary and review cards` — Tasks 6-7
7. `feat(tts): add i18n keys for TTS` — Task 8
