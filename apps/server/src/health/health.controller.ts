import { Controller, Get } from '@nestjs/common';
import { HealthService } from './health.service';
import { HealthResDto } from './dto/health.dto';
import { ApiEnvelope } from 'src/common/swagger/envelope/api-envelope.decorator';
import { Public } from 'src/auth/public.decorator';

@Controller('health')
export class HealthController {
  constructor(private readonly healthService: HealthService) {}

  @Public()
  @ApiEnvelope(HealthResDto, { description: 'Service Health' })
  @Get('/')
  health(): Promise<HealthResDto> {
    return this.healthService.health();
  }
}
