# Taskforge

A self-hosted workflow and automation engine designed for engineers, developers, and technical users.

## Quick Start

```bash
# 1. Clone and install dependencies
git clone https://github.com/taskforge/taskforge.git
cd taskforge
pnpm install

# 2. Start infrastructure (PostgreSQL, Redis)
cd deploy/compose && docker compose up -d

# 3. Generate Prisma client and run migrations
cd apps/server && pnpm prisma:generate && pnpm prisma:migrate dev

# 4. Build shared packages
cd ../..
cd packages/db-access && pnpm build
cd ../queue-config && pnpm build
cd ../contracts && pnpm build

# 5. Start the server
cd ../apps/server && pnpm run start:dev
```

The server will be available at `http://localhost:3000`.

## Documentation

- [Architecture Overview](./docs/Architecture.md)
- [Workflow Engine Concepts](./docs/Taskforge%20-%20Workflow%20Engine.md)
- [Queues and Workers Plan](./docs/Taskforge%20-%20Queues%20and%20Workers%20Plan.md)

## Project Structure

```
taskforge/
├── apps/
│   ├── server/           # API server (NestJS + Fastify)
│   └── worker/           # BullMQ worker for executing steps
├── packages/
│   ├── contracts/        # Zod schemas (API contracts, workflow definitions)
│   ├── db-access/        # Prisma service + repositories
│   └── queue-config/     # BullMQ configuration (Redis connection)
├── docs/                 # Design documents and planning notes
├── deploy/               # Docker Compose configurations
└── scripts/              # Utility scripts
```

## Requirements

- **Node.js** 20+
- **pnpm** 9+
- **PostgreSQL** 14+
- **Redis** 7+