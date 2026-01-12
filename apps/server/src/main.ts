import { cleanupOpenApiDoc } from 'nestjs-zod';
import { NestFactory } from '@nestjs/core';
import {
  NestFastifyApplication,
  FastifyAdapter,
} from '@nestjs/platform-fastify';
import { SwaggerModule } from '@nestjs/swagger';
import helmet from '@fastify/helmet';
import { Logger } from 'nestjs-pino';

import { AppModule } from './app.module';
import { SWAGGER_ENDPOINT, config as SwaggerConfig } from './bootstrap/swagger';
import { ConfigService } from '@nestjs/config';
import { Env } from './config/env';
import { ResponseInterceptor } from './common/interceptors/response/response.interceptor';
import { AllExceptionsFilter } from './common/interceptors/all-exceptions/all-exceptions.filter';

const methods = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'];

async function bootstrap() {
  const app = await NestFactory.create<NestFastifyApplication>(
    AppModule,
    new FastifyAdapter(),
    { cors: { methods }, bufferLogs: true },
  );
  const configService: ConfigService<Env> = app.get(ConfigService);

  const version: string = configService.get('VERSION')!;

  app.setGlobalPrefix(`v${version}/api`);
  await app.register(helmet);
  app.useLogger(app.get(Logger));
  app.enableShutdownHooks();
  app.useGlobalInterceptors(new ResponseInterceptor());
  app.useGlobalFilters(new AllExceptionsFilter());

  SwaggerConfig.info.version = version;

  const openApiDoc = SwaggerModule.createDocument(app, SwaggerConfig);

  const documentFactory = () =>
    SwaggerModule.createDocument(app, cleanupOpenApiDoc(openApiDoc));
  SwaggerModule.setup(SWAGGER_ENDPOINT, app, documentFactory);

  await app.listen(configService.get('PORT'));
}

bootstrap().catch((err) => {
  console.error(err);
  process.exit(1);
});
