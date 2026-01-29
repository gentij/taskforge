import { z } from 'zod';

export const HttpMethodSchema = z.enum([
  'GET',
  'POST',
  'PUT',
  'PATCH',
  'DELETE',
]);

export const HttpRequestSpecSchema = z.object({
  method: HttpMethodSchema,
  url: z
    .string()
    .url()
    .refine((value) => value.startsWith('http://') || value.startsWith('https://'), {
      message: 'url must be an absolute http(s) URL',
    }),
  headers: z.record(z.string(), z.string()).optional(),
  query: z
    .record(z.string(), z.union([z.string(), z.number(), z.boolean()]))
    .optional(),
  body: z.unknown().optional(),
  timeoutMs: z.number().int().positive().optional(),
});

export const HttpStepDefinitionSchema = z.object({
  key: z.string().min(1),
  type: z.literal('http'),
  request: HttpRequestSpecSchema,
});

export const StepDefinitionSchema = z.discriminatedUnion('type', [
  HttpStepDefinitionSchema,
]);

export const WorkflowDefinitionSchema = z.object({
  steps: z.array(StepDefinitionSchema).default([]),
});

export type HttpMethod = z.infer<typeof HttpMethodSchema>;
export type HttpRequestSpec = z.infer<typeof HttpRequestSpecSchema>;
export type HttpStepDefinition = z.infer<typeof HttpStepDefinitionSchema>;
export type StepDefinition = z.infer<typeof StepDefinitionSchema>;
export type WorkflowDefinition = z.infer<typeof WorkflowDefinitionSchema>;
