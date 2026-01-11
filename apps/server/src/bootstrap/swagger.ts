import { DocumentBuilder } from '@nestjs/swagger';

export const config = new DocumentBuilder()
  .setTitle('Taskforge')
  .setDescription('The Taskforge API description')
  .setVersion('1.0')
  .addTag('Taskforge')
  .build();

export const SWAGGER_ENDPOINT = '/api';
