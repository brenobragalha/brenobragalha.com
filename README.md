# brenobragalha.com

Personal site for [brenobragalha.com](https://brenobragalha.com), plus the small
Go generator that builds it.

## Build

```bash
go run ./cmd/site
go test ./...
```

Preview:

```bash
go run ./cmd/site && python3 -m http.server -d public 8080
```

## Content

Site copy lives in [`content/site.yaml`](content/site.yaml). Essays are
`content/writing/<slug>.md`; the filename is the URL slug, so `hello-world.md`
is published at `/writing/hello-world/`.

```markdown
---
title: "Example essay title"
date: 2026-07-29
description: Optional; falls back to the title.
---

Body in Markdown.
```

Markdown is goldmark with GFM. Images go in `static/images/writing/<slug>/…` and
are referenced as `/images/writing/<slug>/…`. The home page lists the five most
recent essays; `/writing/` lists all of them.

## Deploy

Push to `main` (or `workflow_dispatch`) publishes via GitHub Actions. The custom
domain comes from `static/CNAME`.
