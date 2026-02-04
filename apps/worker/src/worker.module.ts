import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { BullModule } from '@nestjs/bullmq';
import { STEP_RUN_QUEUE_NAME, QueueConfigModule } from '@taskforge/queue-config';
import { PrismaModule } from './prisma/prisma.module';
import {
  StepRunRepository,
  WorkflowRunRepository,
  WorkflowVersionRepository,
} from '@taskforge/db-access';
import { ExecutorRegistry } from './executors/executor-registry';
import { HttpExecutorModule } from './executors/http/http-executor.module';
import { TransformExecutorModule } from './executors/transform/transform-executor.module';
import { StepRunProcessor } from './processors/step-run.processor';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: '.env',
    }),
    PrismaModule,
    QueueConfigModule,
    BullModule.registerQueue({ name: STEP_RUN_QUEUE_NAME }),
    HttpExecutorModule,
    TransformExecutorModule,
  ],
  providers: [
    StepRunProcessor,
    StepRunRepository,
    WorkflowRunRepository,
    WorkflowVersionRepository,
    ExecutorRegistry,
  ],
  exports: [],
})
export class WorkerModule {}
