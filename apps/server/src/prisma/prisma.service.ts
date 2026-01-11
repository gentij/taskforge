import {
  Injectable,
  OnModuleInit,
  OnModuleDestroy,
  Logger,
} from '@nestjs/common';
import { PrismaClient } from '@prisma/client';
import { PrismaPg } from '@prisma/adapter-pg';
import { DbHealthDto } from './dto/prisma.dto';

@Injectable()
export class PrismaService
  extends PrismaClient
  implements OnModuleInit, OnModuleDestroy
{
  private readonly logger = new Logger(PrismaService.name);

  constructor() {
    const adapter = new PrismaPg({
      connectionString: process.env.DATABASE_URL!,
    });

    super({ adapter });
  }

  async healthInfo(): Promise<DbHealthDto> {
    const start = Date.now();
    try {
      await this.$queryRaw`SELECT 1`;
      return { ok: true, latencyMs: Date.now() - start };
    } catch (e) {
      return {
        ok: false,
        latencyMs: Date.now() - start,
        error: e instanceof Error ? e.message : String(e),
      };
    }
  }

  async onModuleInit() {
    await this.$connect();

    this.logger.log('Connected to PostgreSQL via Prisma');
  }

  async onModuleDestroy() {
    await this.$disconnect();

    this.logger.log('Disconnected from PostgreSQL');
  }
}
