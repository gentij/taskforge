import { NestFactory } from '@nestjs/core';
import { Logger } from '@nestjs/common';

import { AppModule } from './app.module';

async function bootstrap() {
  const app = await NestFactory.createApplicationContext(AppModule, {
    bufferLogs: true,
  });

  const logger = new Logger('WorkerBootstrap');
  logger.log('Taskforge Worker started');

  app.enableShutdownHooks();
}

bootstrap().catch((err) => {
  console.error(err);
  process.exit(1);
});
