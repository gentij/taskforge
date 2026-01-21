import { DocumentBuilder } from '@nestjs/swagger';

export const config = new DocumentBuilder()
  .setTitle('Taskforge')
  .setDescription('The Taskforge API description')
  .addTag('Taskforge')
  .addBearerAuth(
    {
      type: 'http',
      scheme: 'bearer',
      bearerFormat: 'API Token',
      description: 'Paste your Taskforge API token here',
    },
    'bearer', // <- name of the security scheme
  )
  .build();

export const SWAGGER_ENDPOINT = '/api';
