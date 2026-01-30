import { NestFactory } from '@nestjs/core';
import { Logger } from 'nestjs-pino';

import { AppModule } from './app.module';

async function bootstrap() {
  const app = await NestFactory.createApplicationContext(AppModule, {
    bufferLogs: true,
  });

  const logger = app.get(Logger);
  app.useLogger(logger);

  logger.log('Taskforge Worker started');

  app.enableShutdownHooks();
}

bootstrap().catch((err) => {
  console.error(err);
  process.exit(1);
});