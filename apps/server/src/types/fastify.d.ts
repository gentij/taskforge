import type { ApiToken } from '@prisma/client';

declare module 'fastify' {
  interface FastifyRequest {
    apiToken?: ApiToken;
  }
}
