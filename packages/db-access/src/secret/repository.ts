import { Injectable } from '@nestjs/common';
import type { Prisma, Secret } from '@prisma/client';
import { PrismaService } from '../prisma.service';

@Injectable()
export class SecretRepository {
  constructor(private readonly prisma: PrismaService) {}

  create(data: Prisma.SecretCreateInput): Promise<Secret> {
    return this.prisma.secret.create({ data });
  }

  findMany(): Promise<Secret[]> {
    return this.prisma.secret.findMany({ orderBy: { createdAt: 'desc' } });
  }

  findById(id: string): Promise<Secret | null> {
    return this.prisma.secret.findUnique({ where: { id } });
  }

  findByName(name: string): Promise<Secret | null> {
    return this.prisma.secret.findUnique({ where: { name } });
  }

  findManyByNames(names: string[]): Promise<Secret[]> {
    if (names.length === 0) return Promise.resolve([]);
    return this.prisma.secret.findMany({ where: { name: { in: names } } });
  }

  update(id: string, data: Prisma.SecretUpdateInput): Promise<Secret> {
    return this.prisma.secret.update({ where: { id }, data });
  }

  delete(id: string): Promise<Secret> {
    return this.prisma.secret.delete({ where: { id } });
  }
}
