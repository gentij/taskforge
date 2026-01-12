import { ApiProperty } from '@nestjs/swagger';
import { DbHealthDto } from 'src/prisma/dto/prisma.dto';

export class HealthResDto {
  @ApiProperty({ example: 'ok', enum: ['ok', 'degraded'] })
  status: 'ok' | 'degraded';

  @ApiProperty({ example: '0.1.0' })
  version: string;

  @ApiProperty({ example: 123.45 })
  uptime: number;

  @ApiProperty({ type: DbHealthDto })
  db: DbHealthDto;
}
