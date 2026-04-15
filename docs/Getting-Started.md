---
id: getting-started
title: Getting Started
description: Install Lune and run your first workflow in minutes.
slug: /getting-started
---

# Getting Started

This guide is for operators and users who want to install Lune and run workflows.

If you want to contribute to Lune itself, use `docs/Development.md`.

## 10-Minute Path

If `lune` is already installed, this is the fastest path from zero to first run.

```bash
# 1) Initialize and start stack (external datastores)
lune init \
  --database-url "postgresql://user:pass@db.example.com:5432/lune" \
  --redis-url "redis://redis.example.com:6379"
lune status

# 2) Verify auth works
lune auth whoami

# 3) Create a minimal workflow definition
cat > /tmp/lune-definition.json <<'JSON'
{
  "input": {
    "apiBase": "https://jsonplaceholder.typicode.com"
  },
  "steps": [
    {
      "key": "fetch_post",
      "type": "http",
      "request": {
        "method": "GET",
        "url": "{{input.apiBase}}/posts/1"
      }
    }
  ]
}
JSON

# 4) Create workflow and run it
lune workflow create --name "First Workflow" --definition /tmp/lune-definition.json
lune workflow list
# copy workflow ID from list output
lune workflow run <workflow-id>

# 5) Inspect execution
lune run list <workflow-id>
lune step list <workflow-id> <run-id>
```

## What You Need

- Docker + Docker Compose
- A Lune CLI binary

## Install the CLI

Choose one method.

### Homebrew

```bash
brew tap gentij/lune
brew install lune
```

### AUR (Arch Linux)

```bash
yay -S lune-cli-bin
```

### GitHub Release Binary

Download the matching release artifact from:

- `https://github.com/gentij/lune/releases`

Extract `lune` and put it on your `PATH`.

## Initialize Stack

```bash
lune init \
  --database-url "postgresql://user:pass@db.example.com:5432/lune" \
  --redis-url "redis://redis.example.com:6379"
```

This creates local config and starts the Lune stack (server, worker).

If you want Lune to run bundled local Postgres and Redis instead, use:

```bash
lune init --with-local-datastores
```

Check status:

```bash
lune status
```

Verify auth:

```bash
lune auth whoami
```

## Run Your First Workflow

Create a minimal definition:

```bash
cat > /tmp/lune-definition.json <<'JSON'
{
  "input": {
    "apiBase": "https://jsonplaceholder.typicode.com"
  },
  "steps": [
    {
      "key": "fetch_post",
      "type": "http",
      "request": {
        "method": "GET",
        "url": "{{input.apiBase}}/posts/1"
      }
    }
  ]
}
JSON
```

Create and run:

```bash
lune workflow create --name "First Workflow" --definition /tmp/lune-definition.json
lune workflow list
# use workflow id from list output
lune workflow run <workflow-id>
```

Inspect run details:

```bash
lune run list <workflow-id>
lune step list <workflow-id> <run-id>
```

## Optional: Configure Public Webhook Ingress

Create a webhook trigger and rotate a path key:

```bash
lune trigger create <workflow-id> --type WEBHOOK --name "Inbound"
lune trigger webhook rotate-key <workflow-id> <trigger-id>
```

This prints a URL in the format:

- `http://<host>:3000/v1/api/hooks/:workflowId/:triggerId/:webhookKey`

Use that URL in your webhook provider. For reverse proxy starters, see:

- `deploy/ingress/nginx.lune.conf.example`
- `deploy/ingress/Caddyfile.example`

## Optional: Open the TUI

```bash
lune tui
```

## Next Steps

- [CLI Usage](./CLI-Usage.md)
- [Workflow Definitions](./Lune%20-%20Workflow%20Definitions.md)
- [TUI Guide](./Lune%20-%20TUI.md)
- Full CLI reference: `apps/cli/README.md`
