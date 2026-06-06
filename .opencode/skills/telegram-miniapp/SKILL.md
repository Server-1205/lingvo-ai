---
name: telegram-miniapp
description: Use when working with Telegram Mini App SDK, initData authentication, WebApp events, theme integration, or Telegram-specific UI patterns. NOT for general React development or backend auth logic.
---

# Telegram Mini App Development

## SDK Setup

```typescript
import { init, backButton, themeParams, initData, mainButton } from '@telegram-apps/sdk'

init()

// Access theme colors
const bg = themeParams.state()?.bgColor
const text = themeParams.state()?.textColor
const btn = themeParams.state()?.buttonColor

// Get user data
const user = initData.state()?.user
```

## Key Patterns

### Theme-aware components
```tsx
function ThemedView({ children }: { children: React.ReactNode }) {
  const theme = themeParams.state()
  return (
    <div style={{
      backgroundColor: theme?.bgColor || '#fff',
      color: theme?.textColor || '#000',
      minHeight: '100vh',
    }}>
      {children}
    </div>
  )
}
```

### Native back button
```tsx
import { backButton } from '@telegram-apps/sdk'

function useBackButton(onBack: () => void) {
  useEffect(() => {
    backButton.mount()
    backButton.show()
    const off = backButton.onClick(onBack)
    return () => { off(); backButton.hide() }
  }, [onBack])
}
```

### Sending initData to backend
Include `X-Telegram-Init-Data` header with `initData.raw()` on every API request.

### Viewport
Mini App viewport is typically 390px wide. Design mobile-first. Use `viewportStableHeight` from SDK for proper height.
