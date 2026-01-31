import { z } from 'zod';

export const HttpMethodSchema = z.enum(['GET', 'POST', 'PUT', 'PATCH', 'DELETE']);

export const HttpRequestSpecSchema = z.object({
  method: HttpMethodSchema,
  url: z
    .string()
    .url()
    .refine((value) => value.startsWith('http://') || value.startsWith('https://'), {
      message: 'url must be an absolute http(s) URL',
    }),
  headers: z.record(z.string(), z.string()).optional(),
  query: z.record(z.string(), z.union([z.string(), z.number(), z.boolean()])).optional(),
  body: z.unknown().optional(),
  timeoutMs: z.number().int().positive().optional(),
});

export type HttpMethod = z.infer<typeof HttpMethodSchema>;
export type HttpRequestSpec = z.infer<typeof HttpRequestSpecSchema>;

export const HttpExecutorInputSchema = z.object({
  request: HttpRequestSpecSchema,
  input: z.unknown().default({}),
});

export type HttpExecutorInput = z.infer<typeof HttpExecutorInputSchema>;
