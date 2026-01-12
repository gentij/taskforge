import { Injectable } from '@nestjs/common';
import { HealthResDto } from './dto/health.dto';
import { ConfigService } from '@nestjs/config';
import { PrismaService } from 'src/prisma/prisma.service';

@Injectable()
export class HealthService {
  constructor(
    private readonly configService: ConfigService,
    private readonly prisma: PrismaService,
  ) {}

  async health(): Promise<HealthResDto> {
    const dbOk = await this.prisma.healthInfo();

    return {
      status: 'ok',
      uptime: process.uptime(),
      version: this.configService.get('VERSION', '0.1.0'),
      db: dbOk,
    };
  }
}
