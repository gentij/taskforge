import { Module } from '@nestjs/common';
import { BullModule } from '@nestjs/bullmq';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { Env } from 'src/config/env';
import { STEP_RUN_QUEUE_NAME } from './queue.constants';
import { StepRunQueueService } from './step-run-queue.service';

@Module({
  imports: [
    BullModule.forRootAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: (configService: ConfigService<Env>) => {
        const redisUrl =
          configService.getOrThrow<Env['REDIS_URL']>('REDIS_URL');
        const url = new URL(redisUrl);
        const port = url.port ? Number(url.port) : 6379;

        return {
          connection: {
            host: url.hostname,
            port,
            username: url.username || undefined,
            password: url.password || undefined,
            tls: url.protocol === 'rediss:' ? {} : undefined,
          },
        };
      },
    }),
    BullModule.registerQueue({ name: STEP_RUN_QUEUE_NAME }),
  ],
  providers: [StepRunQueueService],
  exports: [StepRunQueueService],
})
export class QueueModule {}
