-- +goose Up
CREATE TABLE channels (
    channel_id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE videos (
    video_id TEXT PRIMARY KEY,
    channel_id TEXT NOT NULL REFERENCES channels(channel_id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    is_summarized BOOLEAN NOT NULL DEFAULT FALSE,
    published_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- NEW: The summaries table Claude mentioned
CREATE TABLE summaries (
    video_id TEXT PRIMARY KEY REFERENCES videos(video_id) ON DELETE CASCADE,
    summary_text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_digested BOOLEAN NOT NULL DEFAULT FALSE 
);

-- +goose Down
DROP TABLE summaries;
DROP TABLE videos;
DROP TABLE channels;