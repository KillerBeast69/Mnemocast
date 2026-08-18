-- name: CreateChannel :exec
INSERT INTO channels (channel_id, title)
VALUES ($1, $2)
ON CONFLICT (channel_id) DO NOTHING;

-- name: CreateVideo :exec
INSERT INTO videos (video_id, channel_id, title, url)
VALUES ($1, $2, $3, $4)
ON CONFLICT (video_id) DO NOTHING;

-- name: CreateSummary :exec
INSERT INTO summaries (video_id, summary_text)
VALUES ($1, $2)
ON CONFLICT (video_id) DO NOTHING;