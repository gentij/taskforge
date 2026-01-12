export type ErrorDetail = { field?: string; message: string };

export type HttpExceptionResponseObject = {
  message?: unknown;
  code?: unknown;
  details?: unknown;
};
