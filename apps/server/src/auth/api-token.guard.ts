import { CanActivate, ExecutionContext, Injectable } from '@nestjs/common';
import { Reflector } from '@nestjs/core';
import type { FastifyRequest } from 'fastify';

import { ApiTokenService } from 'src/api-token/api-token.service';
import { AppError } from 'src/common/http/errors/app-error';
import { ErrorDefinitions } from 'src/common/http/errors/error-codes';
import { CryptoService } from 'src/crypto/crypto.service';

export const IS_PUBLIC_KEY = 'isPublic';

@Injectable()
export class ApiTokenGuard implements CanActivate {
  constructor(
    private readonly reflector: Reflector,
    private readonly apiTokenService: ApiTokenService,
    private readonly cryptoService: CryptoService,
  ) {}

  async canActivate(context: ExecutionContext): Promise<boolean> {
    const isPublic = this.reflector.getAllAndOverride<boolean>(IS_PUBLIC_KEY, [
      context.getHandler(),
      context.getClass(),
    ]);
    if (isPublic) return true;

    const req = context.switchToHttp().getRequest<FastifyRequest>();

    const rawToken = this.extractBearer(req);
    if (!rawToken) {
      throw AppError.unauthorized(ErrorDefinitions.AUTH.MISSING_BEARER_TOKEN);
    }

    const tokenHash = this.cryptoService.hashApiToken(rawToken);
    const apiToken = await this.apiTokenService.validateTokenHash(tokenHash);

    if (!apiToken) {
      throw AppError.unauthorized(ErrorDefinitions.AUTH.INVALID_TOKEN);
    }

    req.apiToken = apiToken;
    return true;
  }

  private extractBearer(req: FastifyRequest): string | null {
    const header = req.headers.authorization;
    if (typeof header !== 'string') return null;

    const [scheme, value] = header.split(' ');
    if (scheme?.toLowerCase() !== 'bearer') return null;
    if (!value) return null;

    return value.trim();
  }
}
