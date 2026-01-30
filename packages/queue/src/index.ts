import { Module } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { BullModule } from '@nestjs/bullmq';

export const STEP_RUN_QUEUE_NAME = 'step-runs';

export const QueueConfigModule = BullModule.forRootAsync({
  imports: [ConfigModule],
  inject: [ConfigService],
  useFactory: (configService: ConfigService) => {
    const redisUrl = configService.getOrThrow<string>('REDIS_URL');
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
});

export const registerStepRunQueue = () =>
  BullModule.registerQueue({ name: STEP_RUN_QUEUE_NAME });