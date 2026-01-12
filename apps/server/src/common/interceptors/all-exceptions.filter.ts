import {
  ArgumentsHost,
  Catch,
  ExceptionFilter,
  HttpException,
  HttpStatus,
  Logger,
} from '@nestjs/common';
import type { FastifyReply, FastifyRequest } from 'fastify';
import { ErrorDetail } from './all-exceptions.types';
import {
  isHttpExceptionResponseObject,
  isObjectRecord,
  isStringArray,
  toStringSafe,
} from 'src/lib/utils/util';

@Catch()
export class AllExceptionsFilter implements ExceptionFilter {
  private readonly logger = new Logger(AllExceptionsFilter.name);

  catch(exception: unknown, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const req = ctx.getRequest<FastifyRequest>();
    const res = ctx.getResponse<FastifyReply>();

    let statusCode = HttpStatus.INTERNAL_SERVER_ERROR;
    let code = 'INTERNAL_ERROR';
    let message = 'Something went wrong';
    let details: ErrorDetail[] | undefined;

    if (exception instanceof HttpException) {
      statusCode = exception.getStatus();

      const response = exception.getResponse(); // unknown-ish
      if (isHttpExceptionResponseObject(response)) {
        // ValidationPipe: { message: string[], error: 'Bad Request', statusCode: 400 }
        if (isStringArray(response.message)) {
          code = 'VALIDATION_ERROR';
          message = 'Request validation failed';
          details = response.message.map((m) => ({ message: m }));
        } else {
          if (typeof response.code === 'string') code = response.code;
          if (typeof response.message === 'string') message = response.message;
          if (Array.isArray(response.details)) {
            // best-effort normalize details
            details = response.details
              .map((d): ErrorDetail | null => {
                if (!isObjectRecord(d)) return null;
                const msg =
                  typeof d.message === 'string' ? d.message : undefined;
                if (!msg) return null;
                const field = typeof d.field === 'string' ? d.field : undefined;
                return { field, message: msg };
              })
              .filter((x): x is ErrorDetail => x !== null);
          }
        }
      } else if (typeof response === 'string') {
        // Sometimes getResponse() is a string
        code = 'HTTP_ERROR';
        message = response;
      } else {
        code = 'HTTP_ERROR';
        message = exception.message;
      }
    } else if (exception instanceof Error) {
      // Don’t leak internals to client, but do log them.
      this.logger.error(exception.message, exception.stack);
      // keep generic message/code for client
    } else {
      this.logger.error(`Non-Error thrown: ${toStringSafe(exception)}`);
    }

    res.status(statusCode).send({
      ok: false as const,
      statusCode,
      path: req.url,
      timestamp: new Date().toISOString(),
      error: {
        code,
        message,
        ...(details ? { details } : {}),
      },
      meta: { requestId: req.id },
    });
  }
}
