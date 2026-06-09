-- internal/db/schema.sql
-- Lingvo AI database schema (SQLite)
-- Executed on every server start. Idempotent.

CREATE TABLE IF NOT EXISTS users (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    telegram_id   INTEGER NOT NULL UNIQUE,
    username      TEXT DEFAULT '',
    lang          TEXT NOT NULL DEFAULT 'uz' CHECK(lang IN ('uz','ru')),
    level         TEXT NOT NULL DEFAULT 'a1' CHECK(level IN ('a1','a2','b1','b2','c1')),
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id       INTEGER NOT NULL REFERENCES users(id),
    date          TEXT NOT NULL,
    count         INTEGER NOT NULL DEFAULT 0,
    UNIQUE(user_id, date)
);

CREATE TABLE IF NOT EXISTS subscriptions (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id       INTEGER NOT NULL UNIQUE REFERENCES users(id),
    plan          TEXT NOT NULL CHECK(plan IN ('weekly','monthly')),
    stars_amount  INTEGER NOT NULL,
    started_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at    DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS vocabulary (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id       INTEGER NOT NULL REFERENCES users(id),
    word          TEXT NOT NULL,
    translation   TEXT NOT NULL,
    example       TEXT NOT NULL,
    level         TEXT DEFAULT 'a1',
    review_count  INTEGER DEFAULT 0,
    ease_factor   REAL DEFAULT 2.5,
    interval      INTEGER DEFAULT 0,
    last_reviewed_at DATETIME,
    next_review   DATETIME,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, word)
);

CREATE TABLE IF NOT EXISTS daily_progress (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id       INTEGER NOT NULL REFERENCES users(id),
    date          TEXT NOT NULL,
    messages_sent INTEGER DEFAULT 0,
    words_learned INTEGER DEFAULT 0,
    quizzes_taken INTEGER DEFAULT 0,
    UNIQUE(user_id, date)
);

CREATE INDEX IF NOT EXISTS idx_messages_user_date ON messages(user_id, date);
CREATE INDEX IF NOT EXISTS idx_vocab_user ON vocabulary(user_id);
CREATE INDEX IF NOT EXISTS idx_progress_user_date ON daily_progress(user_id, date);

CREATE TABLE IF NOT EXISTS error_history (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER NOT NULL REFERENCES users(id),
    original        TEXT NOT NULL,
    corrected       TEXT NOT NULL,
    category        TEXT NOT NULL DEFAULT 'grammar',
    severity        TEXT NOT NULL DEFAULT 'minor',
    rule_violated   TEXT DEFAULT '',
    learning_tip    TEXT DEFAULT '',
    context         TEXT DEFAULT '',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_error_user_date ON error_history(user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_error_category ON error_history(user_id, category);

ALTER TABLE vocabulary ADD COLUMN ease_factor REAL DEFAULT 2.5;
ALTER TABLE vocabulary ADD COLUMN interval INTEGER DEFAULT 0;
ALTER TABLE vocabulary ADD COLUMN last_reviewed_at DATETIME;

CREATE TABLE IF NOT EXISTS ielts_scores (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER NOT NULL REFERENCES users(id),
    module          TEXT NOT NULL CHECK(module IN ('writing_task1','writing_task2','speaking','reading')),
    band_score      REAL NOT NULL,
    details         TEXT DEFAULT '{}',
    prompt          TEXT DEFAULT '',
    user_response   TEXT DEFAULT '',
    feedback        TEXT DEFAULT '',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_ielts_user_module ON ielts_scores(user_id, module);
CREATE INDEX IF NOT EXISTS idx_ielts_user_date ON ielts_scores(user_id, created_at);
