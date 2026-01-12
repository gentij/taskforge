import {
  ApiBadRequestResponse,
  ApiConflictResponse,
  ApiForbiddenResponse,
  ApiInternalServerErrorResponse,
  ApiNotFoundResponse,
  ApiTooManyRequestsResponse,
  ApiUnauthorizedResponse,
} from '@nestjs/swagger';
import { RespDecoratorFactory } from './api-envelope.types';

export type EnvelopeError = 400 | 401 | 403 | 404 | 409 | 429 | 500;

export const errorDecorators: Record<EnvelopeError, RespDecoratorFactory> = {
  400: ApiBadRequestResponse,
  401: ApiUnauthorizedResponse,
  403: ApiForbiddenResponse,
  404: ApiNotFoundResponse,
  409: ApiConflictResponse,
  429: ApiTooManyRequestsResponse,
  500: ApiInternalServerErrorResponse,
};

export const errorExamples: Record<EnvelopeError, unknown> = {
  400: {
    ok: false,
    statusCode: 400,
    example: '/{path}',
    timestamp: new Date().toISOString(),
    error: {
      code: 'VALIDATION_ERROR',
      message: 'Request validation failed',
      details: [{ field: 'name', message: 'name must be a string' }],
    },
    meta: { requestId: 'req-123' },
  },
  401: {
    ok: false,
    statusCode: 401,
    example: '/{path}',
    timestamp: new Date().toISOString(),
    error: { code: 'UNAUTHORIZED', message: 'Missing or invalid token' },
    meta: { requestId: 'req-123' },
  },
  403: {
    ok: false,
    statusCode: 403,
    example: '/{path}',
    timestamp: new Date().toISOString(),
    error: {
      code: 'FORBIDDEN',
      message: 'You do not have access to this resource',
    },
    meta: { requestId: 'req-123' },
  },
  404: {
    ok: false,
    statusCode: 404,
    example: '/{path}',
    timestamp: new Date().toISOString(),
    error: { code: 'NOT_FOUND', message: 'Resource not found' },
    meta: { requestId: 'req-123' },
  },
  409: {
    ok: false,
    statusCode: 409,
    example: '/{path}',
    timestamp: new Date().toISOString(),
    error: { code: 'CONFLICT', message: 'Unique constraint violation' },
    meta: { requestId: 'req-123' },
  },
  429: {
    ok: false,
    statusCode: 429,
    example: '/{path}',
    timestamp: new Date().toISOString(),
    error: { code: 'RATE_LIMITED', message: 'Too many requests' },
    meta: { requestId: 'req-123' },
  },
  500: {
    ok: false,
    statusCode: 500,
    example: '/{path}',
    timestamp: new Date().toISOString(),
    error: { code: 'INTERNAL_ERROR', message: 'Something went wrong' },
    meta: { requestId: 'req-123' },
  },
};
