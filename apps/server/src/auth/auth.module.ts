import { Module } from '@nestjs/common';
import { APP_GUARD } from '@nestjs/core';
import { ApiTokenGuard } from './api-token.guard';
import { ApiTokenModule } from 'src/api-token/api-token.module';
import { CryptoModule } from 'src/crypto/crypto.module';
import { AuthController } from './auth.controller';
import { AuthService } from './auth.service';

@Module({
  imports: [ApiTokenModule, CryptoModule],
  controllers: [AuthController],
  providers: [AuthService, { provide: APP_GUARD, useClass: ApiTokenGuard }],
})
export class AuthModule {}
