import { Module } from '@nestjs/common';
import { PrismaModule } from 'src/prisma/prisma.module';
import { WorkflowModule } from 'src/workflow/workflow.module';
import { OrchestrationModule } from 'src/core/orchestration.module';
import { TriggerController } from './trigger.controller';
import { TriggerRepository } from '@taskforge/db-access';
import { TriggerService } from './trigger.service';
import { CronTriggerScheduler } from './cron/cron-trigger.scheduler';

@Module({
  imports: [PrismaModule, WorkflowModule, OrchestrationModule],
  controllers: [TriggerController],
  providers: [TriggerService, TriggerRepository, CronTriggerScheduler],
  exports: [TriggerService, TriggerRepository],
})
export class TriggerModule {}
