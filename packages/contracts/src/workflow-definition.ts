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
    .refine((value) => {
      // Allow template patterns like {{input.apiUrl}}/path
      if (value.includes('{{') && value.includes('}}')) {
        return true;
      }
      // Must be valid URL
      try {
        new URL(value);
        return value.startsWith('http://') || value.startsWith('https://');
      } catch {
        return false;
      }
    }, {
      message: 'url must be an absolute http(s) URL or a template pattern',
    }),
  headers: z.record(z.string(), z.string()).optional(),
  query: z
    .record(z.string(), z.union([z.string(), z.number(), z.boolean()]))
    .optional(),
  body: z.unknown().optional(),
  timeoutMs: z.number().int().positive().optional(),
});

export const BaseStepDefinitionSchema = z.object({
  key: z
    .string()
    .min(1)
    .regex(/^[a-zA-Z0-9_-]+$/, {
      message: 'key must contain only letters, numbers, underscore, or hyphen',
    }),
  dependsOn: z.array(z.string()).optional(),
  input: z.record(z.string(), z.unknown()).optional(),
});

export const HttpStepDefinitionSchema = BaseStepDefinitionSchema.extend({
  type: z.literal('http'),
  request: HttpRequestSpecSchema,
});

export const TransformRequestSpecSchema = z
  .object({
    source: z.record(z.string(), z.unknown()).optional(),
    output: z.unknown(),
  })
  .strict();

export const TransformStepDefinitionSchema = BaseStepDefinitionSchema.extend({
  type: z.literal('transform'),
  request: TransformRequestSpecSchema,
});

export const StepDefinitionSchema = z.discriminatedUnion('type', [
  HttpStepDefinitionSchema,
  TransformStepDefinitionSchema,
]);

export const WorkflowDefinitionSchema = z.object({
  input: z.record(z.string(), z.unknown()).optional(),
  steps: z.array(StepDefinitionSchema).default([]),
});

export type HttpMethod = z.infer<typeof HttpMethodSchema>;
export type HttpRequestSpec = z.infer<typeof HttpRequestSpecSchema>;
export type HttpStepDefinition = z.infer<typeof HttpStepDefinitionSchema>;
export type TransformRequestSpec = z.infer<typeof TransformRequestSpecSchema>;
export type TransformStepDefinition = z.infer<typeof TransformStepDefinitionSchema>;
export type StepDefinition = z.infer<typeof StepDefinitionSchema>;
export type WorkflowDefinition = z.infer<typeof WorkflowDefinitionSchema>;
