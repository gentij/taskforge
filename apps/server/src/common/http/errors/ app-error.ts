import { HttpException, HttpStatus } from '@nestjs/common';
import { AppErrorDetails, ErrorDefinition } from './error.types';
import { ErrorDefinitions } from './error-codes';

type AppErrorCtorOpts = {
  statusCode: number;
  error: ErrorDefinition;
  details?: AppErrorDetails;
};

export class AppError extends HttpException {
  readonly code: ErrorDefinition['code'];

  constructor(opts: AppErrorCtorOpts) {
    super(
      {
        code: opts.error.code,
        message: opts.error.message,
        details: opts.details,
      },
      opts.statusCode,
    );

    this.code = opts.error.code;
  }

  static badRequest(
    error: ErrorDefinition = ErrorDefinitions.COMMON.BAD_REQUEST,
    details?: AppErrorDetails,
  ) {
    return new AppError({
      statusCode: HttpStatus.BAD_REQUEST,
      error,
      details,
    });
  }

  static unauthorized(details?: AppErrorDetails) {
    return new AppError({
      statusCode: HttpStatus.UNAUTHORIZED,
      error: ErrorDefinitions.COMMON.UNAUTHORIZED,
      details,
    });
  }

  static forbidden(details?: AppErrorDetails) {
    return new AppError({
      statusCode: HttpStatus.FORBIDDEN,
      error: ErrorDefinitions.COMMON.FORBIDDEN,
      details,
    });
  }

  static notFound(
    error: ErrorDefinition = ErrorDefinitions.COMMON.NOT_FOUND,
    details?: AppErrorDetails,
  ) {
    return new AppError({
      statusCode: HttpStatus.NOT_FOUND,
      error,
      details,
    });
  }

  static conflict(error: ErrorDefinition, details?: AppErrorDetails) {
    return new AppError({
      statusCode: HttpStatus.CONFLICT,
      error,
      details,
    });
  }

  static tooManyRequests(details?: AppErrorDetails) {
    return new AppError({
      statusCode: HttpStatus.TOO_MANY_REQUESTS,
      error: ErrorDefinitions.COMMON.RATE_LIMITED,
      details,
    });
  }

  static internal(details?: AppErrorDetails) {
    return new AppError({
      statusCode: HttpStatus.INTERNAL_SERVER_ERROR,
      error: ErrorDefinitions.COMMON.INTERNAL_ERROR,
      details,
    });
  }
}
