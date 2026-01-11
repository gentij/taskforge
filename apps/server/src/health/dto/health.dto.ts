import { ApiProperty } from '@nestjs/swagger';

export class HealthResDto {
  @ApiProperty()
  status: 'ok' | 'degraded';

  @ApiProperty()
  version: string;

  @ApiProperty()
  uptime: number;

  @ApiProperty()
  timestamp: string;
}
