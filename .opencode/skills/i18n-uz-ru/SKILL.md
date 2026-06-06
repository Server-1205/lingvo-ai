---
name: i18n-uz-ru
description: Use when adding, modifying, or reviewing translation keys in uz.json and ru.json for Lingvo AI. Ensures both languages are always in sync and follows i18next conventions with react-i18next. NOT for general React component work or backend logic.
---

# Internationalization (Uzbek + Russian)

## File Locations
- `web/src/locales/uz.json` — Uzbek (Latin script)
- `web/src/locales/ru.json` — Russian

## Key Structure (flat keys, namespaced by screen)

```json
{
  "chat.title": "AI English Tutor",
  "chat.input_placeholder": "Type a message...",
  "chat.send": "Send",
  "chat.corrections": "Corrections",
  "chat.limit_reached": "Daily limit reached. Upgrade for unlimited access.",
  "nav.chat": "Chat",
  "nav.vocabulary": "Vocabulary",
  "nav.progress": "Progress",
  "nav.profile": "Profile",
  "vocab.title": "My Words",
  "vocab.search": "Search...",
  "vocab.add": "Add word",
  "progress.title": "My Progress",
  "progress.level": "Level",
  "progress.streak": "Day streak",
  "subscription.title": "Premium",
  "subscription.weekly": "Weekly",
  "subscription.monthly": "Monthly",
  "subscription.subscribe": "Subscribe",
  "common.loading": "Loading...",
  "common.error": "Something went wrong",
  "common.retry": "Retry"
}
```

## Usage in React

```tsx
import { useTranslation } from 'react-i18next'

function ChatScreen() {
  const { t } = useTranslation()

  return (
    <div>
      <h1>{t('chat.title')}</h1>
      <input placeholder={t('chat.input_placeholder')} />
      <button>{t('chat.send')}</button>
    </div>
  )
}
```

## Rules
1. **Every new key** must be added to BOTH uz.json and ru.json
2. Keys are kebab-case, dot-separated by screen
3. Use `{{variable}}` syntax for interpolation: `"chat.corrections_count": "{{count}} corrections found"`
4. Plurals: use i18next plural syntax: `"correction": "correction", "correction_plural": "corrections"`
5. Uzbek = Latin script (not Cyrillic)
6. Default language: `uz`
