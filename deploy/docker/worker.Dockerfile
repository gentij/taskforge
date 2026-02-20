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
RUN pnpm -C apps/worker build

FROM node:20-slim AS runner

WORKDIR /app
ENV NODE_ENV=production

COPY --from=builder /app/node_modules /app/node_modules
COPY --from=builder /app/packages /app/packages
COPY --from=builder /app/apps/worker/node_modules /app/apps/worker/node_modules
COPY --from=builder /app/apps/worker/dist /app/apps/worker/dist
COPY --from=builder /app/apps/worker/package.json /app/apps/worker/package.json

CMD ["node", "apps/worker/dist/src/main.js"]
