# Mnemocast

> *mnemo- (memory) + -cast (broadcast) — so you actually remember what these channels said, without watching all of it.*

**Status: 📋 Planning — no code written yet.** This README exists to pin down scope and architecture before implementation starts, so the plan doesn't quietly grow once building begins.

## The problem

I'm subscribed to channels I actually want to keep up with — the kind that cover what's happening right now in tools, languages, and the industry (ThePrimeagen, Theo, Fireship, that tier). I don't have time to watch all of it, titles and thumbnails don't tell me enough to know if a video matters, and by the time I'd get to it, the moment's usually passed.

Mnemocast watches the channels for me, reads the transcript, and mails me a summary short enough to read over coffee.

## What it does (v1 scope)

1. Give it a list of channels to track.
2. On a schedule, it checks each channel for new uploads.
3. For each new video, it pulls the transcript.
4. An LLM turns the transcript into a short, plain-language summary.
5. Once a day (or week), everything new gets bundled into one email.

No video should take longer to "catch up on" than the length of its summary.

## Roadmap

### v1 — MVP
- [ ] Manually configured channel list (no OAuth yet)
- [ ] Scheduled job polls for each channel's uploads via a cheap endpoint (not `search.list`)
- [ ] Transcript fetch per new video
- [ ] LLM summarization (~150 words per video)
- [ ] Daily digest email, simple HTML template

### v2 — Stop polling
- [ ] Move from polling to YouTube's PubSubHubbub/WebSub push notifications
- [ ] Subscription renewal job — leases expire roughly every 10 days and need re-subscribing
- [ ] Quota-aware request layer: cache aggressively, avoid expensive endpoints (`search.list` costs 100 units against a 10,000/day default budget)
- [ ] Idempotency so a video never gets summarized twice if a notification fires more than once

### v3 — Make it actually mine
- [ ] Sign in with Google, auto-import real subscriptions instead of a manual list
- [ ] Per-channel preferences (mute a channel, change frequency)
- [ ] Web archive of past digests, not just email
- [ ] Per-channel digest cadence (some daily, some weekly)

## Architecture (planned)

```
YouTube upload
      │
      ▼
┌──────────────────────────┐
│ Poller (v1) /             │
│ Webhook receiver (v2)     │
└────────────┬──────────────┘
             │ new video event
             ▼
      ┌──────────────┐
      │  Job queue    │
      └──────┬───────┘
             ▼
┌──────────────────────────┐
│ Worker                    │
│  - fetch transcript       │
│  - summarize (LLM)        │
│  - store result           │
└────────────┬──────────────┘
             ▼
      ┌──────────────┐
      │  Postgres     │  (channels, videos, summaries, users)
      └──────┬───────┘
             ▼
┌──────────────────────────┐
│ Digest scheduler           │
│ (daily / weekly cron)      │
└────────────┬──────────────┘
             ▼
      ┌──────────────┐
      │ Email sender  │
      └──────────────┘
```

## Tech stack (planned, not final)

| Layer | Choice | Why |
|---|---|---|
| Language | Go or Python | Whichever more comfortable shipping in |
| Database | PostgreSQL | Channels/videos/summaries/users — relational fits fine |
| Queue | Redis + a worker, or a DB-backed queue | Don't need Kafka for one person's subscriptions |
| Scheduling | Cron (v1) → PubSubHubbub webhook (v2) | Start simple, earn the complexity later |
| Email | SES, Postmark, or SendGrid | Any transactional email provider |
| Packaging | Docker | Matches the backend path |
| CI/CD | GitHub Actions | Same |

## Project layout (planned)

```
mnemocast/
├── ingest/        # YouTube API client, transcript fetch, poller/webhook handler
├── summarize/     # LLM prompt + summarization logic
├── digest/        # batching + email templating
├── store/         # Postgres access layer, migrations
├── docker-compose.yml
└── README.md
```

## Known constraints going in

- **YouTube Data API quota**: 10,000 units/day by default, and `search.list` alone costs 100 units per call — the polling design has to avoid it.
- **PubSubHubbub leases are short-lived** (~10 days) — v2 needs a renewal job, not a one-time subscribe.
- **Transcripts miss anything visual** — code on screen, diagrams, charts won't make it into the summary. Worth stating that limitation in the digest itself rather than overselling it.

## Status / next steps

- [ ] Repo scaffolding + Docker Compose (Postgres + app)
- [ ] DB schema (channels, videos, summaries)
- [ ] v1 end to end: poll → transcript → summarize → one email
- [ ] v2: webhook migration + quota-aware caching
- [ ] v3: OAuth + preferences

## License

