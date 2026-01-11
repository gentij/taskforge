import { ApiProperty } from '@nestjs/swagger';

export class DbHealthDto {
  @ApiProperty({ example: true })
  ok: boolean;

  @ApiProperty({ example: 12 })
  latencyMs: number;

  @ApiProperty({
    required: false,
    example: 'connection refused',
    nullable: true,
  })
  error?: string;
}
