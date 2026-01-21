import { Injectable } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';
import { ApiToken } from '@prisma/client';
import { CreateApiTokenInput } from './api-token.types';

@Injectable()
export class ApiTokenRepository {
  constructor(private readonly prisma: PrismaService) {}

  findAll(): Promise<ApiToken[]> {
    return this.prisma.apiToken.findMany();
  }

  findActive(): Promise<ApiToken[]> {
    return this.prisma.apiToken.findMany({
      where: { revokedAt: null },
    });
  }

  findByHash(tokenHash: string): Promise<ApiToken | null> {
    return this.prisma.apiToken.findUnique({
      where: { tokenHash },
    });
  }

  create({ name, tokenHash, scopes }: CreateApiTokenInput): Promise<ApiToken> {
    return this.prisma.apiToken.create({
      data: {
        name,
        tokenHash,
        scopes: scopes ?? [],
      },
    });
  }

  updateLastUsed(id: string): Promise<ApiToken> {
    return this.prisma.apiToken.update({
      where: { id },
      data: { lastUsedAt: new Date() },
    });
  }

  revoke(id: string): Promise<ApiToken> {
    return this.prisma.apiToken.update({
      where: { id },
      data: { revokedAt: new Date() },
    });
  }
}
