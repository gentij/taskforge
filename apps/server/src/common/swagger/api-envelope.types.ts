import { ApiBadRequestResponse } from '@nestjs/swagger';

export type EnvelopeError = 400 | 401 | 403 | 404 | 409 | 429 | 500;
export type RespDecoratorFactory = (
  options?: Parameters<typeof ApiBadRequestResponse>[0],
) => MethodDecorator;
