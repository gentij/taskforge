import { HttpExceptionResponseObject } from 'src/common/interceptors/all-exceptions/all-exceptions.types';

export function isObjectRecord(v: unknown): v is Record<string, unknown> {
  return typeof v === 'object' && v !== null;
}

export function isHttpExceptionResponseObject(
  v: unknown,
): v is HttpExceptionResponseObject {
  return isObjectRecord(v);
}

export function isStringArray(v: unknown): v is string[] {
  return Array.isArray(v) && v.every((x) => typeof x === 'string');
}

export function toStringSafe(v: unknown): string {
  return typeof v === 'string' ? v : 'Unknown error';
}
