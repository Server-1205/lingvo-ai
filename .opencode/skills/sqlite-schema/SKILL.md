---
name: sqlite-schema
description: Use when creating or modifying SQLite schema, adding tables, columns, indexes, or writing SQL queries for the Lingvo AI project. NOT for Go handler logic or application code.
---

# SQLite Schema & Migrations

## Schema File (`internal/db/schema.sql`)

```sql
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    telegram_id INTEGER UNIQUE NOT NULL,
    username TEXT,
    lang TEXT DEFAULT 'uz' CHECK(lang IN ('uz', 'ru')),
    level TEXT DEFAULT 'a1' CHECK(level IN ('a1','a2','b1','b2','c1')),
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);
```

## Common Queries (sqlx)

```go
// Named query with struct
type User struct {
    TelegramID int64  `db:"telegram_id"`
    Lang       string `db:"lang"`
    Level      string `db:"level"`
}

user := User{TelegramID: 12345}
err := db.Get(&user, `SELECT * FROM users WHERE telegram_id = :telegram_id`, user)

// Insert with named params
_, err := db.NamedExec(`
    INSERT INTO users (telegram_id, username, lang)
    VALUES (:telegram_id, :username, :lang)
    ON CONFLICT(telegram_id) DO UPDATE SET
        username = excluded.username,
        updated_at = datetime('now')
`, user)
```

## Migration Patterns

Add column:
```sql
ALTER TABLE users ADD COLUMN timezone TEXT DEFAULT 'Asia/Tashkent';
```

New table: append to the end of schema.sql. The `Migrate()` function executes all statements.

## Rules
- Always use `CREATE TABLE IF NOT EXISTS`
- Always use `CHECK` constraints for enums
- Index columns used in WHERE clauses
- Use TEXT for dates (ISO 8601 format: `datetime('now')`)
- Foreign keys: `REFERENCES users(id) ON DELETE CASCADE`
