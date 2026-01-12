import { ApiProperty } from '@nestjs/swagger';

export class ApiErrorItemDto {
  @ApiProperty({ required: false, example: 'name' })
  field?: string;

  @ApiProperty({ example: 'name must be a string' })
  message: string;
}

export class ApiErrorDto {
  @ApiProperty()
  code: string;

  @ApiProperty({ example: 'Something went wrong' })
  message: string;

  @ApiProperty({ required: false, type: [ApiErrorItemDto] })
  details?: ApiErrorItemDto[];
}

export class ApiErrorMetaDto {
  @ApiProperty({ required: false, example: 'req-123' })
  requestId?: string;
}

export class ApiErrorResponseDto {
  @ApiProperty({ example: false })
  ok: false;

  @ApiProperty({ description: 'HTTP status code for the error' })
  statusCode: number;

  @ApiProperty({ example: '/{path}' })
  path: string;

  @ApiProperty()
  timestamp: string;

  @ApiProperty({ type: ApiErrorDto })
  error: ApiErrorDto;

  @ApiProperty({ type: ApiErrorMetaDto, required: false })
  meta?: ApiErrorMetaDto;
}
