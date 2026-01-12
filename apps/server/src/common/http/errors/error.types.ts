export type FieldDetail = { field?: string; message: string };

export type AppErrorDetails = FieldDetail[] | Record<string, unknown> | string;

export type ErrorDefinition = {
  code: string;
  message: string;
};
