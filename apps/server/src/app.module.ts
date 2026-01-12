import { APP_PIPE, APP_INTERCEPTOR } from '@nestjs/core';
import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { ZodValidationPipe, ZodSerializerInterceptor } from 'nestjs-zod';
import { LoggerModule } from 'nestjs-pino';

import { validateEnv } from './config/env';
import { HealthModule } from './health/health.module';
import { CoreModule } from './core/core.module';
import { PrismaModule } from './prisma/prisma.module';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      validate: validateEnv,
    }),
    LoggerModule.forRoot(),
    CoreModule,
    PrismaModule,
    HealthModule,
  ],
  providers: [
    {
      provide: APP_PIPE,
      useClass: ZodValidationPipe,
    },
    {
      provide: APP_INTERCEPTOR,
      useClass: ZodSerializerInterceptor,
    },
  ],
})
export class AppModule {}
