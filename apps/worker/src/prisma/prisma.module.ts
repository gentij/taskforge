import { Global, Module } from '@nestjs/common';
import { PrismaService } from '@taskforge/db-access';

@Global()
@Module({
  providers: [PrismaService],
  exports: [PrismaService],
})
export class PrismaModule {}
