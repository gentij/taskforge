import { ApiProperty } from '@nestjs/swagger';

export class ApiMetaDto {
  @ApiProperty({ required: false, example: 'req-123' })
  requestId?: string;
}

export class ApiResponseDto<T> {
  @ApiProperty({ example: true })
  ok: boolean;

  @ApiProperty({ example: 200 })
  statusCode: number;

  @ApiProperty({ example: '/{path}' })
  path: string;

  @ApiProperty({ example: '2026-01-12T22:30:00.000Z' })
  timestamp: string;

  // override this with ApiOkResponse schema
  data: T;

  @ApiProperty({ type: ApiMetaDto })
  meta: ApiMetaDto;
}
