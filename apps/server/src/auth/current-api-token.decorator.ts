import { createParamDecorator, ExecutionContext } from '@nestjs/common';
import type { FastifyRequest } from 'fastify';
import type { ApiToken } from '@prisma/client';

export const CurrentApiToken = createParamDecorator(
  (_: unknown, ctx: ExecutionContext): ApiToken | undefined => {
    const req = ctx.switchToHttp().getRequest<FastifyRequest>();
    return req.apiToken;
  },
);
