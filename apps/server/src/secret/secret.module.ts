import { Module } from '@nestjs/common';
import { PrismaModule } from 'src/prisma/prisma.module';
import { SecretController } from './secret.controller';
import { SecretRepository } from '@taskforge/db-access';
import { SecretService } from './secret.service';

@Module({
  imports: [PrismaModule],
  controllers: [SecretController],
  providers: [SecretService, SecretRepository],
  exports: [SecretService],
})
export class SecretModule {}
