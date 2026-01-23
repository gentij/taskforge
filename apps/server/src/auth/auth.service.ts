import { Injectable } from '@nestjs/common';
import { ApiToken } from '@prisma/client';
import { WhoamiResDto } from './dto/auth.dto';

@Injectable()
export class AuthService {
  constructor() {}

  whoami(token: ApiToken): WhoamiResDto {
    return { id: token.id, name: token.name, scopes: token.scopes };
  }
}
