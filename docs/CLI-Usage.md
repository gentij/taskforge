---
id: cli-usage
title: CLI Usage
description: Common Lune CLI commands for day-to-day operations.
slug: /cli
---

# CLI Usage

This page is a practical command reference for running Lune.

For full command coverage, see `apps/cli/README.md`.

## Global Flags

- `--server`: API base URL (default `http://localhost:3000/v1/api`)
- `--output`: `table` or `json`
- `--quiet`: print IDs only
- `--no-color`: disable colored output
- `--config`: config file path

## Pagination and Sorting

List commands support pagination and server-side sorting.

Common list flags:

- `--page`
- `--page-size`
- `--sort-by`
- `--sort-order` (`asc|desc`)

Supported `--sort-by` values:

- `workflow list`: `createdAt|updatedAt`
- `workflow version list`: `version|createdAt`
- `trigger list`: `createdAt|updatedAt`
- `run list`: `createdAt|updatedAt`
- `step list`: `createdAt|updatedAt`
- `secret list`: `createdAt|updatedAt`

Examples:

```bash
lune workflow list --sort-by updatedAt --sort-order desc
lune run list <workflow-id> --sort-by createdAt --sort-order asc
lune workflow version list <workflow-id> --sort-by version --sort-order desc
```

## Stack Commands

Run once (external datastores):

```bash
lune init \
  --database-url "postgresql://user:pass@db.example.com:5432/lune" \
  --redis-url "redis://redis.example.com:6379"
```

Local bundled Postgres/Redis (optional):

```bash
lune init --with-local-datastores
```

Day-to-day:

```bash
lune start
lune status
lune logs --follow
lune stop
```

## Auth Commands

```bash
lune auth login --token "<LUNE_ADMIN_TOKEN>"
lune auth whoami
lune auth status
lune auth logout
```

## Workflow Lifecycle

```bash
lune workflow create --name "My Workflow" --definition definition.json
lune workflow list
lune workflow get <workflow-id>
lune workflow update <workflow-id> --name "New Name"
lune workflow run <workflow-id> --input input.json
lune workflow delete <workflow-id>
```

Notes:

- `--input` values override colliding keys from workflow definition `input`.
- `--overrides` is for per-step HTTP request overrides (`query`/`body`) keyed by step key.

## Run Notifications

Workflow definitions support optional `notifications` entries for run completion events.

- Providers: `discord`, `slack`
- Events: `SUCCEEDED`, `FAILED`
- `webhook` can be absolute `http(s)` or `{{secret.NAME}}`

Use `lune workflow create` or `lune workflow version create` with a definition that includes `notifications`.

## Triggers

```bash
lune trigger create <workflow-id> --type CRON --name "Nightly" --config cron.json
lune trigger create <workflow-id> --type WEBHOOK --name "Inbound"
lune trigger webhook rotate-key <workflow-id> <trigger-id>
lune trigger list <workflow-id>
lune trigger get <workflow-id> <trigger-id>
lune trigger update <workflow-id> <trigger-id> --is-active=false
lune trigger delete <workflow-id> <trigger-id>
```

Public webhook ingress path format:

- `POST /v1/api/hooks/:workflowId/:triggerId/:webhookKey`

After rotating a webhook key, call the generated URL directly from your webhook provider.

## Runs and Steps

```bash
lune run list <workflow-id>
lune run get <workflow-id> <run-id>
lune step list <workflow-id> <run-id>
lune step get <workflow-id> <run-id> <step-run-id>
```

## Secrets

```bash
lune secret create --name API_KEY --value "secret"
lune secret list
lune secret get <secret-id>
lune secret update <secret-id> --description "Rotated"
lune secret delete <secret-id>
```

## TUI

```bash
lune tui
```
