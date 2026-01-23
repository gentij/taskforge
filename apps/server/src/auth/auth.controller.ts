import { Controller, Get } from '@nestjs/common';
import { ApiEnvelope } from 'src/common/swagger/envelope/api-envelope.decorator';
import { AuthService } from './auth.service';
import { CurrentApiToken } from './current-api-token.decorator';
import type { ApiToken } from '@prisma/client';
import { WhoamiResDto } from './dto/auth.dto';
import { ApiBearerAuth } from '@nestjs/swagger';

@Controller('auth')
@ApiBearerAuth('bearer')
export class AuthController {
  constructor(private readonly authService: AuthService) {}

  @ApiEnvelope(WhoamiResDto, { description: 'WHOMAI', errors: [500, 401] })
  @Get('/whoami')
  whoami(@CurrentApiToken() token: ApiToken): WhoamiResDto {
    return this.authService.whoami(token);
  }
}
