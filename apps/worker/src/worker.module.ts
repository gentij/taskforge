import { Module } from '@nestjs/common';
import { BullModule } from '@nestjs/bullmq';
import { STEP_RUN_QUEUE_NAME, QueueConfigModule } from '@taskforge/queue-config';

@Module({
  imports: [
    QueueConfigModule,
    BullModule.registerQueue({ name: STEP_RUN_QUEUE_NAME }),
  ],
  controllers: [],
  providers: [],
  exports: [],
})
export class WorkerModule {}