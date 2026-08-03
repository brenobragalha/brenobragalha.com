# brenobragalha.com

Personal site and the small Go generator that builds it.

## Build

```bash
go run ./cmd/site
go test ./...
```

Output goes to `public/` (gitignored). Preview:

```bash
go run ./cmd/site && python3 -m http.server -d public 8080
```

## Tests

One end-to-end fixture build: checks that the site generates (pages, feed, Markdown) and republishes cleanly. It does not cover validation failures—those fail at `go run ./cmd/site`.

## Content

Site copy: [`content/site.yaml`](content/site.yaml) — requires `base_url`, `name`, and
`contact` (`github`, `x`, `linkedin`, `email`).

Essays: `content/writing/<slug>.md`. Filename is the URL slug (lowercase kebab-case).
Front matter needs `title` and `date` (`YYYY-MM-DD`); `description` is optional and
falls back to the title.

```markdown
---
title: "Example essay title"
date: 2026-07-29
description: Optional; falls back to the title.
---

Body in Markdown.
```

Images: `static/images/writing/<slug>/…` → `/images/writing/<slug>/…`.
Home lists the five most recent essays; `/writing/` lists all.

## Deploy

Push to `main` publishes via GitHub Actions. Domain: `static/CNAME`.
