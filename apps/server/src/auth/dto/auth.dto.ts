import { ApiProperty } from '@nestjs/swagger';

export class WhoamiResDto {
  @ApiProperty({ example: '' })
  id: string;

  @ApiProperty({ example: '' })
  name: string;

  @ApiProperty({ example: [''] })
  scopes: string[];
}
