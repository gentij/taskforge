import { HttpStatus } from '@nestjs/common';
import { Prisma } from '@prisma/client';
import { ErrorDefinitions } from './error-codes';

function isPrismaKnownRequestError(
  e: unknown,
): e is Prisma.PrismaClientKnownRequestError {
  return e instanceof Prisma.PrismaClientKnownRequestError;
}

function isPrismaValidationError(
  e: unknown,
): e is Prisma.PrismaClientValidationError {
  return e instanceof Prisma.PrismaClientValidationError;
}

function isPrismaInitError(
  e: unknown,
): e is Prisma.PrismaClientInitializationError {
  return e instanceof Prisma.PrismaClientInitializationError;
}

function isPrismaRustPanicError(
  e: unknown,
): e is Prisma.PrismaClientRustPanicError {
  return e instanceof Prisma.PrismaClientRustPanicError;
}

export function mapPrismaError(e: unknown): {
  statusCode: number;
  code: string;
  message: string;
  details?: unknown;
} | null {
  if (isPrismaKnownRequestError(e)) {
    // Prisma codes: https://www.prisma.io/docs/reference/api-reference/error-reference
    switch (e.code) {
      case 'P2002': {
        // Unique constraint failed
        return {
          statusCode: HttpStatus.CONFLICT,
          code: ErrorDefinitions.DATABASE.UNIQUE_CONSTRAINT.code,
          message: ErrorDefinitions.DATABASE.UNIQUE_CONSTRAINT.message,
          details: { prismaCode: e.code, meta: e.meta },
        };
      }

      case 'P2025': {
        // Record not found (e.g. update/delete where no record)
        return {
          statusCode: HttpStatus.NOT_FOUND,
          code: ErrorDefinitions.COMMON.NOT_FOUND.code,
          message: ErrorDefinitions.COMMON.NOT_FOUND.message,
          details: { prismaCode: e.code, meta: e.meta },
        };
      }

      case 'P2003': {
        // Foreign key constraint failed
        return {
          statusCode: HttpStatus.CONFLICT,
          code: ErrorDefinitions.COMMON.CONFLICT?.code ?? 'CONFLICT',
          message:
            ErrorDefinitions.COMMON.CONFLICT?.message ??
            'Foreign key constraint failed',
          details: { prismaCode: e.code, meta: e.meta },
        };
      }

      default:
        return {
          statusCode: HttpStatus.BAD_REQUEST,
          code: ErrorDefinitions.COMMON.BAD_REQUEST.code,
          message: ErrorDefinitions.COMMON.BAD_REQUEST.message,
          details: { prismaCode: e.code, meta: e.meta },
        };
    }
  }

  if (isPrismaValidationError(e)) {
    return {
      statusCode: HttpStatus.BAD_REQUEST,
      code: ErrorDefinitions.COMMON.BAD_REQUEST.code,
      message: ErrorDefinitions.COMMON.BAD_REQUEST.message,
      details: { prisma: 'VALIDATION', message: e.message },
    };
  }

  if (isPrismaInitError(e)) {
    // DB connection / credentials / schema engine init etc.
    return {
      statusCode: HttpStatus.SERVICE_UNAVAILABLE,
      code: ErrorDefinitions.COMMON.INTERNAL_ERROR.code,
      message: 'Database unavailable',
      details: { prisma: 'INIT', message: e.message },
    };
  }

  if (isPrismaRustPanicError(e)) {
    return {
      statusCode: HttpStatus.INTERNAL_SERVER_ERROR,
      code: ErrorDefinitions.COMMON.INTERNAL_ERROR.code,
      message: ErrorDefinitions.COMMON.INTERNAL_ERROR.message,
      details: { prisma: 'RUST_PANIC' },
    };
  }

  return null;
}
