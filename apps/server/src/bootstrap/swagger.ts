import { DocumentBuilder } from '@nestjs/swagger';

export const config = new DocumentBuilder()
  .setTitle('Taskforge')
  .setDescription('The Taskforge API description')
  .addTag('Taskforge')
  .addBearerAuth()
  .build();

export const SWAGGER_ENDPOINT = '/api';
