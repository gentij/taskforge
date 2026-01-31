import { Module } from '@nestjs/common';

import { PrismaModule } from 'src/prisma/prisma.module';
import { QueueModule } from 'src/queue/queue.module';

import { OrchestrationService } from './orchestration.service';

@Module({
  imports: [PrismaModule, QueueModule],
  providers: [OrchestrationService],
  exports: [OrchestrationService],
})
export class OrchestrationModule {}
