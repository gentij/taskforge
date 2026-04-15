import { DocumentBuilder } from '@nestjs/swagger';

export const config = new DocumentBuilder()
  .setTitle('Lune')
  .setDescription('The Lune API description')
  .addTag('Lune')
  .addBearerAuth(
    {
      type: 'http',
      scheme: 'bearer',
      bearerFormat: 'API Token',
      description: 'Paste your Lune API token here',
    },
    'bearer', // <- name of the security scheme
  )
  .build();

export const SWAGGER_ENDPOINT = '/api';
