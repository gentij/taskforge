import { Module } from '@nestjs/common';
import { ApiTokenService } from './api-token.service';
import { ApiTokenRepository } from './api-token.repository';

@Module({
  providers: [ApiTokenService, ApiTokenRepository],
  exports: [ApiTokenService],
})
export class ApiTokenModule {}
