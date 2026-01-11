import { Controller, Get } from '@nestjs/common';
import { HealthService } from './health.service';
import { HealthResDto } from './dto/health.dto';
import { ApiOkResponse } from '@nestjs/swagger';

@Controller('health')
export class HealthController {
  constructor(private readonly healthService: HealthService) {}

  @ApiOkResponse({ type: HealthResDto })
  @Get('/')
  health(): HealthResDto {
    return this.healthService.health();
  }
}
