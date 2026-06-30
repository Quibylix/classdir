CREATE TABLE IF NOT EXISTS schema_migrations (
    version TEXT PRIMARY KEY
);

CREATE TABLE presentations (
    id         UUID NOT NULL PRIMARY KEY,
    title      TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE slides (
    id              UUID NOT NULL PRIMARY KEY,
    presentation_id UUID NOT NULL REFERENCES presentations(id) ON DELETE CASCADE,
    slide_number    INTEGER NOT NULL,
    content         TEXT NOT NULL,
    metadata        JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(presentation_id, slide_number)
);

CREATE TABLE students (
    id              UUID NOT NULL PRIMARY KEY,
    presentation_id UUID NOT NULL REFERENCES presentations(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(presentation_id, name)
);
