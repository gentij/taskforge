FROM node:20-slim AS builder

WORKDIR /app

RUN corepack enable

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY packages ./packages
COPY apps ./apps

RUN pnpm install --frozen-lockfile
RUN pnpm -C apps/server prisma:generate
RUN pnpm -C packages/db-access build \
  && pnpm -C packages/queue build \
  && pnpm -C packages/contracts build
RUN pnpm -C apps/server build

FROM node:20-slim AS runner

WORKDIR /app
ENV NODE_ENV=production

RUN apt-get update \
  && apt-get install -y openssl \
  && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/node_modules /app/node_modules
COPY --from=builder /app/packages /app/packages
COPY --from=builder /app/apps/server/node_modules /app/apps/server/node_modules
COPY --from=builder /app/apps/server/prisma /app/apps/server/prisma
COPY --from=builder /app/apps/server/prisma.config.ts /app/apps/server/prisma.config.ts
COPY --from=builder /app/apps/server/dist /app/apps/server/dist
COPY --from=builder /app/apps/server/package.json /app/apps/server/package.json

EXPOSE 3000

CMD ["node", "apps/server/dist/src/main.js"]
