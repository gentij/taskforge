import {
  Injectable,
  OnModuleInit,
  OnModuleDestroy,
  Logger,
} from '@nestjs/common';
import { PrismaClient } from '@prisma/client';
import { PrismaPg } from '@prisma/adapter-pg';
import { DbHealthDto } from './dto/prisma.dto';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class PrismaService
  extends PrismaClient
  implements OnModuleInit, OnModuleDestroy
{
  private readonly logger = new Logger(PrismaService.name);

  constructor(private readonly configService: ConfigService) {
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
    const db = await this.healthInfo();

    if (db.ok) {
      this.logger.log(`Connected to PostgreSQL (latency ${db.latencyMs}ms)`);
      return;
    }

    this.logger.error(
      `PostgreSQL NOT reachable: ${db.error ?? 'unknown error'}`,
    );
  }

  async onModuleDestroy() {
    await this.$disconnect();

    this.logger.log('Disconnected from PostgreSQL');
  }
}
