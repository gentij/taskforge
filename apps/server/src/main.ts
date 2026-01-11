import { cleanupOpenApiDoc } from 'nestjs-zod';
import { NestFactory } from '@nestjs/core';
import {
  NestFastifyApplication,
  FastifyAdapter,
} from '@nestjs/platform-fastify';
import { SwaggerModule } from '@nestjs/swagger';
import helmet from '@fastify/helmet';

import { AppModule } from './app.module';
import { SWAGGER_ENDPOINT, config as SwaggerConfig } from './bootstrap/swagger';
import { ConfigService } from '@nestjs/config';
import { Env } from './config/env';

const methods = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'];

async function bootstrap() {
  const app = await NestFactory.create<NestFastifyApplication>(
    AppModule,
    new FastifyAdapter(),
    { cors: { methods } },
  );

  app.setGlobalPrefix('v1/api/');
  await app.register(helmet);

  const configService: ConfigService<Env> = app.get(ConfigService);

  SwaggerConfig.info.version = configService.get('VERSION')!;

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
