# Lune CLI

The Lune CLI lets you manage workflows, triggers, runs, steps, and secrets from the terminal.

For self-hosted setup, run `lune init` first to provision `~/.lune/` and stack config.

## Quickstart

```bash
lune auth login --token "<LUNE_ADMIN_TOKEN>"
lune workflow list
```

If your API is not on the default port, pass `--server`:

```bash
lune --server http://localhost:3100/v1/api auth whoami
```

Create a workflow from a definition file:

```bash
lune workflow create --name "My Workflow" --definition definition.json
```

Create a CRON trigger:

```bash
lune trigger create <workflow-id> --type CRON --name "Nightly" --config cron.json
```

Run a workflow and check runs/steps:

```bash
lune workflow run <workflow-id> --input input.json
lune run list <workflow-id>
lune step list <workflow-id> <run-id>
```

## Global Flags

- `--output table|json` (default: `table`)
- `--quiet` (print IDs only)
- `--no-color` (disable colored status output)
- `--server` (API base URL, default: `http://localhost:3000/v1/api`)
- `--config` (config file path)

`NO_COLOR=1` disables color output as well.

## Pagination and Sorting

List commands support pagination and server-side sorting.

Common flags:

- `--page` (default: `1`)
- `--page-size` (default: `25`)
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
lune workflow list --page 1 --page-size 25 --sort-by updatedAt --sort-order desc
lune run list <workflow-id> --sort-by createdAt --sort-order asc
lune workflow version list <workflow-id> --sort-by version --sort-order desc
```

## Auth

```bash
lune auth login
lune auth status
lune auth whoami
lune auth logout
```

## Stack

`lune init` must be run once before using stack commands.

By default, init expects external Postgres and Redis URLs:

```bash
lune init \
  --database-url "postgresql://user:pass@db.example.com:5432/lune" \
  --redis-url "redis://redis.example.com:6379"
```

To run bundled local Postgres/Redis containers instead:

```bash
lune init --with-local-datastores
```

```bash
lune start
lune start server worker
lune start --foreground
lune status
lune logs --follow
lune stop
```

## Workflows

```bash
lune workflow list
lune workflow get <workflow-id>
lune workflow create --name "My Workflow" --definition definition.json
lune workflow update <workflow-id> --name "New Name"
lune workflow update <workflow-id> --is-active=false
lune workflow delete <workflow-id>
lune workflow run <workflow-id> --input input.json --overrides overrides.json
lune workflow validate <workflow-id> --definition definition.json
```

Notes:

- `--input` values override colliding keys from workflow definition `input`.
- `--overrides` applies only to `http` steps and supports request `query`/`body` overrides keyed by step key.

Sample `definition.json`:

```json
{
  "input": {
    "apiBase": "https://jsonplaceholder.typicode.com",
    "postId": 1
  },
  "steps": [
    {
      "key": "fetch_post",
      "type": "http",
      "request": {
        "method": "GET",
        "url": "{{input.apiBase}}/posts/{{input.postId}}"
      }
    }
  ]
}
```

Sample with run notifications:

```json
{
  "input": {
    "shouldPass": true
  },
  "notifications": [
    {
      "provider": "discord",
      "webhook": "{{secret.DISCORD_WEBHOOK_URL}}",
      "on": ["FAILED"]
    },
    {
      "provider": "slack",
      "webhook": "{{secret.SLACK_WEBHOOK_URL}}",
      "on": ["SUCCEEDED", "FAILED"]
    }
  ],
  "steps": [
    {
      "key": "gate",
      "type": "condition",
      "request": {
        "expr": "input.shouldPass",
        "assert": true
      }
    }
  ]
}
```

Notes:

- `webhook` accepts absolute `http(s)` URLs or `{{secret.NAME}}` references.
- Notifications are evaluated on final run status (`SUCCEEDED`/`FAILED`).
- Delivery is best-effort and does not change workflow run status.

## Workflow Versions

```bash
lune workflow version list <workflow-id>
lune workflow version get <workflow-id> <version>
lune workflow version create <workflow-id> --definition definition.json
```

## Triggers

```bash
lune trigger list <workflow-id>
lune trigger get <workflow-id> <trigger-id>
lune trigger create <workflow-id> --type CRON --name "Nightly" --config cron.json
lune trigger create <workflow-id> --type WEBHOOK --name "Inbound"
lune trigger webhook rotate-key <workflow-id> <trigger-id>
lune trigger update <workflow-id> <trigger-id> --is-active=false
lune trigger delete <workflow-id> <trigger-id>
```

Sample `cron.json`:

```json
{
  "cron": "0 2 * * *",
  "timezone": "UTC"
}
```

## Runs

```bash
lune run list <workflow-id>
lune run get <workflow-id> <run-id>
```

## Steps

```bash
lune step list <workflow-id> <run-id>
lune step get <workflow-id> <run-id> <step-run-id>
```

## Secrets

```bash
lune secret list
lune secret get <secret-id>
lune secret create --name API_KEY --value "my-secret"
lune secret update <secret-id> --description "Rotated"
lune secret delete <secret-id>
```

Secrets are not printed in table output. Use `--output json` if you need raw JSON.

## Output Modes

- `--output table` shows human-readable tables (default)
- `--output json` prints raw JSON for scripting
- `--quiet` prints only IDs

## Troubleshooting

- **Token not set**: run `lune auth login`
- **Validation errors**: verify JSON files match the server schema (e.g., CRON uses `cron`, not `expression`)
