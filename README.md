# brenobragalha.com

Personal site for [brenobragalha.com](https://brenobragalha.com), plus the small
Go generator that builds it.

## Build

```bash
go run ./cmd/site
go test ./...
```

Output goes to `public/` (gitignored). Preview:

```bash
go run ./cmd/site && python3 -m http.server -d public 8080
```

## Content

Site copy lives in [`content/site.yaml`](content/site.yaml). Essays are
`content/writing/<slug>.md`. The filename is the URL slug and must be lowercase
kebab-case (`hello-world.md` → `/writing/hello-world/`). Front matter needs
`title` and `date` (`YYYY-MM-DD`); `description` is optional and falls back to
the title.

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
