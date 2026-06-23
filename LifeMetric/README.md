# Health-Tech Agent Roster — Gemini CLI Edition

This is the Gemini CLI version of your nine-person virtual team, matching
the names and roles already defined for Claude Code:

| Name | Role | Background |
|---|---|---|
| Sade (Folasade Adewale) | Product Strategist | Nigeria |
| Arjun (Arjun Mehta) | Solution Architect | India |
| Hana (Hana Kobayashi) | Backend Engineer | Japan |
| Mateus (Mateus Oliveira) | Frontend Engineer | Brazil |
| Lukas (Lukas Hoffmann) | Hardware Integration Engineer | Germany |
| Layla (Layla Haddad) | QA / Test Engineer | Lebanon |
| Amani (Amani Wanjiru) | Security & Compliance Reviewer | Kenya |
| Renata (Renata Cruz) | Data Privacy Engineer | Mexico |
| Tane (Tane Ngata) | DevOps / Release Engineer | New Zealand (Māori) |

## How this differs from the Claude Code version

Gemini CLI's subagent system uses Markdown files with YAML frontmatter
instead of Claude Code's agent config format, but the underlying idea is
the same: each agent has its own system prompt, restricted tool access,
and the main Gemini session (you, talking to Gemini) acts as orchestrator,
deciding which specialist to route a task to based on its description.

Each agent here is scoped to the actual MCP servers already connected in
this project (GitHub, AWS, Postgres, Cloudflare) — but only the ones
relevant to their role. None of them have IAM, billing, or destructive
database access; those stay with you.

## Installation

1. Unzip this into your project root so the files land at:
   `~/dev/agentic-setup/.gemini/agents/*.md`
2. Make sure you're running Gemini CLI from that trusted folder
3. Restart Gemini CLI — it should detect the new agents automatically
4. Run `/agents` to confirm they're all registered

## Usage

Just talk to the main Gemini session naturally — it routes to the right
specialist based on what you ask:

- "Sade, draft a PRD for [feature]" or "what should the MVP scope be for X"
- "Arjun, design the data model for tracking [health metric]"
- "Hana, implement the API for [feature]"
- "Mateus, build the screen for [feature]"
- "Lukas, integrate the [device name] sensor data"
- "Layla, write a test plan for [feature]"
- "Amani, security review this pull request"
- "Renata, check what PII this table is storing"
- "Tane, set up the deployment pipeline for [service]"

You can also be explicit: "have Sade and Arjun look at this together" works
the same way it did in the Claude Code version — naming them by name
prompts the right handoff.

## What's deliberately NOT an agent

Same four roles as the Claude Code roster — these need a real, legally
accountable human, not an AI draft:

- Regulatory officer (FDA/CE submissions, formal claims)
- Clinical safety reviewer (does this meet a clinical safety bar)
- Data protection counsel (formal legal privacy compliance)
- You, as final decision-maker and accountable party

Every agent above is instructed to flag — not silently resolve — anything
that touches these four areas.
