---
title: "A Lightweight Playbook for Shipping LLM Experiments"
date: 2025-11-25T20:45:00-08:00
draft: false
description: "The checklist I run through to keep small language model experiments scoped, observable, and shippable."
tags:
  - ai
  - workflow
  - experiments
images:
  - "/og/HelloWorld.png"
---

I have been shipping tiny LLM-powered tools for friends, clients, and my own curiosity projects. After a dozen iterations, a loose routine emerged. This post captures the checklist so I can hand it to future-me (or you) when a new experiment idea lands in the notes app.

![Illustration of the LLM experiment loop showing framing, scaffolding, instrumentation, and retrospectives.](/images/blog/llm-playbook-loop.svg)
*The four-step loop I try to keep sacred whenever a new idea hits the notebook.*

## 1. Frame the experiment like a scientist

- **Hypothesis** — write one sentence about the behavior you expect, including the measurable signal (CTR improves, handoffs drop, etc.).
- **Guardrails** — write the unacceptable outcome ("hallucinated legal advice") so you can detect it quickly.
- **Stopping rule** — note the moment you will call the experiment done (>200 sessions, 5 interviews, etc.).

Keeping the framing in a short README inside the repo prevents scope drift once you touch the keyboard.

## 2. Build the scaffolding first, model second

I force myself to wire the ingestion pipeline, prompt storage, and logging before I tune the prompt. Two tiny helpers speed this up:

1. **Prompt registry** — a plain YAML file with named prompts. Every experiment loads from it so I can diff prompt changes like code.
2. **Dry-run flag** — a CLI switch that swaps the LLM client with a deterministic stub so I can develop the UX without burning tokens.

The moment the scaffolding exists, model iterations feel cheaper because I am not fighting plumbing.

## 3. Instrument like a scientist, not a marketer

Instead of throwing an analytics SDK at the prototype, I log structured events to a local SQLite file: prompt hash, latency, token counts, and the human feedback snippet. The file ships with the repo and migrates with a single command. When something breaks, I can replay the exact conversation because the context window is saved alongside the metadata.

## 4. Decide with evidence, not vibes

Every experiment gets a pre-written retro template:

- What surprised me?
- What felt expensive?
- What would make this 10× cheaper to validate next time?

I fill it out before publishing any results. This keeps me honest about the real blockers (tooling? data access? reviewer time?) instead of defaulting to "need smarter prompt".

## Tooling stack that kept working

| Layer | Tool | Why |
| --- | --- | --- |
| Data collection | Observable notebooks | Quick charts plus markdown for notes |
| Prompt diffing | Git plus `promptlint` | Fails CI when I forget eval coverage |
| Evaluation | Lightweight rubric scores in CSV | Human-in-the-loop without overhead |
| Hosting | Fly.io Machines | Suspend when idle, resume in ~200 ms |

None of these are flashy, but every tool reduces the activation energy for the next batch of experiments.

## What’s next

The next iteration is borrowing more ideas from traditional experimentation frameworks: pre-registration, better sample sizing, and automated alerting when the guardrails trip. Until then, this playbook is the minimum viable ritual that lets me chase weird LLM ideas without derailing the rest of my roadmap.
