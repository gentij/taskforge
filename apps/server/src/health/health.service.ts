import { Injectable } from '@nestjs/common';
import { HealthResDto } from './dto/health.dto';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class HealthService {
  constructor(private readonly configService: ConfigService) {}

  health(): HealthResDto {
    return {
      status: 'ok',
      timestamp: new Date().toISOString(),
      uptime: process.uptime(),
      version: this.configService.get('VERSION', '0.1.0'),
    };
  }
}
