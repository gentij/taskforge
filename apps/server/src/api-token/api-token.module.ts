import { Module } from '@nestjs/common';
import { ApiTokenService } from './api-token.service';
import { ApiTokenRepository } from '@lune/db-access';
import { PrismaModule } from 'src/prisma/prisma.module';

@Module({
  imports: [PrismaModule],
  providers: [ApiTokenService, ApiTokenRepository],
  exports: [ApiTokenService],
})
export class ApiTokenModule {}
