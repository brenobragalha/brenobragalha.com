---
title: "Porting My Knowledge System from Notion to Obsidian"
date: 2025-11-25T21:05:00-08:00
draft: false
description: "Field notes from migrating a five-year Notion archive into a nimble Obsidian graph without losing context."
tags:
  - pkms
  - tooling
  - workflows
images:
  - "/og/notion-to-obsidian.svg"
---

Over the last five years I defaulted to Notion for everything: project trackers, scratchpads, even grocery lists. The structure was beautiful but the friction grew every time I wanted to remix notes or spotlight weak links across projects. So I spent two weekends migrating to Obsidian. Here is the play-by-play.

![Flow diagram of the Notion to Obsidian migration pipeline showing export, taxonomy, automation, and canvas synthesis.](/images/blog/notion-pipeline.svg)
*My four-stage pipeline for keeping the migration sane.*

## 1. Normalize the Notion export

1. Use the official bulk export, but choose **Markdown & CSV** instead of HTML. This preserves inline formatting while keeping relations as CSV attachments.
2. Run a quick script to flatten `/Subpage/` folders into `YYYY/MM/slug.md`. That timestamp prefix becomes essential for sorting inside Obsidian's file explorer.
3. Convert Notion callouts to standard markdown using a simple `sed` pass. Obsidian's callouts are similar, so the visual diff stays minimal.

## 2. Design link taxonomies first

Before importing anything, I wrote two short documents: `01-links.md` (what deserves a wiki link) and `02-tags.md` (what stays a loose hashtag). For me:

- Wiki links are reserved for recurring people, projects, and named concepts.
- Tags capture transient context such as `#meeting`, `#idea`, or `#draft`.

By drawing that line early, I avoided the temptation to link every noun just because Obsidian made it easy.

## 3. Automate progressive enrichment

Instead of trying to refactor 1,200 notes in one sitting, I added three lightweight automations:

- A daily recurring task that surfaces five random notes for review.
- A dataview query that lists orphaned notes with zero backlinks.
- A command palette shortcut that stamps the current note with a `Last touched` metadata block.

This keeps the migration alive without turning into a one-time death march.

## 4. Lean on canvas maps for synthesis

Obsidian Canvas became the missing artifact I wanted in Notion. I now drag key notes into a canvas, sketch the relationships, and link back to the source markdown. The map acts as a project README that stays in sync with the underlying files.

## 5. Publish a migration retro

The last step was capturing the hiccups:

- Notion's databases export poorly, so I screenshot the critical timelines for posterity.
- Templates with complex formulas rarely translate; I rewrote them using Obsidian Templater.
- Mobile capture feels worse today, but the faster desktop flow compensates.

If you are sitting on a mountain of Notion notes, consider Obsidian not as a fresh start but as a refactor. The markdown files are yours, the graph surfaces new connections, and the plugins keep the system adaptable.
