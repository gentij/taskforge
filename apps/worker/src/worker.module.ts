import { Module } from '@nestjs/common';
import { BullModule } from '@nestjs/bullmq';
import { STEP_RUN_QUEUE_NAME, QueueConfigModule } from '@taskforge/queue-config';
import { PrismaModule } from './prisma/prisma.module';
import { StepRunRepository, WorkflowVersionRepository } from '@taskforge/db-access';
import { ExecutorRegistry } from './executors/executor-registry';
import { HttpExecutorModule } from './executors/http/http-executor.module';
import { StepRunProcessor } from './processors/step-run.processor';

@Module({
  imports: [
    PrismaModule,
    QueueConfigModule,
    BullModule.registerQueue({ name: STEP_RUN_QUEUE_NAME }),
    HttpExecutorModule,
  ],
  providers: [StepRunProcessor, StepRunRepository, WorkflowVersionRepository, ExecutorRegistry],
  exports: [],
})
export class WorkerModule {}
