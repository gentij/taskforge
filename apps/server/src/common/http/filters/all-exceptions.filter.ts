import {
  ArgumentsHost,
  Catch,
  ExceptionFilter,
  HttpException,
  HttpStatus,
  Logger,
} from '@nestjs/common';
import type { FastifyReply, FastifyRequest } from 'fastify';
import { ZodValidationException } from 'nestjs-zod';
import { ZodError } from 'zod';
import { ErrorDefinitions } from 'src/common/http/errors/error-codes';
import {
  isHttpExceptionResponseObject,
  isStringArray,
  toStringSafe,
} from 'src/lib/utils/util';
import { mapPrismaError } from '../errors/prisma-errors.lib';

@Catch()
export class AllExceptionsFilter implements ExceptionFilter {
  private readonly logger = new Logger(AllExceptionsFilter.name);

  catch(exception: unknown, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const req = ctx.getRequest<FastifyRequest>();
    const res = ctx.getResponse<FastifyReply>();

    let statusCode = HttpStatus.INTERNAL_SERVER_ERROR;
    let code: string = ErrorDefinitions.COMMON.INTERNAL_ERROR.code;
    let message: string = ErrorDefinitions.COMMON.INTERNAL_ERROR.message;
    let details: unknown;

    if (exception instanceof ZodValidationException) {
      const zodError = exception.getZodError();
      if (zodError instanceof ZodError) {
        this.logger.warn(`Zod validation failed: ${zodError.message}`);
        details = zodError.issues.map((i) => ({
          field: i.path.join('.'),
          message: i.message,
        }));
        code = ErrorDefinitions.COMMON.VALIDATION_ERROR.code;
        message = ErrorDefinitions.COMMON.VALIDATION_ERROR.message;
        statusCode = exception.getStatus?.() ?? HttpStatus.BAD_REQUEST;
      }
    }

    if (exception instanceof HttpException) {
      statusCode = exception.getStatus();
      const response = exception.getResponse();

      if (isHttpExceptionResponseObject(response)) {
        if (isStringArray(response.message)) {
          code = ErrorDefinitions.COMMON.VALIDATION_ERROR.code;
          message = ErrorDefinitions.COMMON.VALIDATION_ERROR.message;
          details = response.message.map((m) => ({ message: m }));
        } else {
          if (typeof response.code === 'string') code = response.code;
          if (typeof response.message === 'string') message = response.message;
          if (typeof response.details !== 'undefined')
            details = response.details;
        }
      } else if (typeof response === 'string') {
        code = ErrorDefinitions.COMMON.BAD_REQUEST.code;
        message = response;
      } else {
        code = ErrorDefinitions.COMMON.BAD_REQUEST.code;
        message = exception.message;
      }
    } else if (exception instanceof Error) {
      this.logger.error(exception.message, exception.stack);
    } else {
      this.logger.error(`Non-Error thrown: ${toStringSafe(exception)}`);
    }

    const prismaMapped = mapPrismaError(exception);

    if (prismaMapped) {
      statusCode = prismaMapped.statusCode;
      code = prismaMapped.code;
      message = prismaMapped.message;
      details = prismaMapped.details;

      res.status(statusCode).send({
        ok: false as const,
        statusCode,
        path: req.url,
        timestamp: new Date().toISOString(),
        error: {
          code,
          message,
          ...(typeof details !== 'undefined' ? { details } : {}),
        },
        meta: { requestId: req.id },
      });
      return;
    }

    res.status(statusCode).send({
      ok: false as const,
      statusCode,
      path: req.url,
      timestamp: new Date().toISOString(),
      error: {
        code,
        message,
        ...(typeof details !== 'undefined' ? { details } : {}),
      },
      meta: { requestId: req.id },
    });
  }
}
