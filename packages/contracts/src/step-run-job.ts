import { z } from 'zod';

export const HttpRequestOverrideSchema = z
  .object({
    query: z
      .record(z.string(), z.union([z.string(), z.number(), z.boolean()]))
      .optional(),
    body: z.unknown().optional(),
  })
  .strict();

export const StepRunJobPayloadSchema = z.object({
  workflowRunId: z.string().min(1),
  stepRunId: z.string().min(1),
  stepKey: z.string().min(1),
  workflowVersionId: z.string().min(1),
  input: z.unknown().default({}),
  dependsOn: z.array(z.string()).default([]),
  requestOverride: HttpRequestOverrideSchema.optional(),
});

export type StepRunJobPayload = z.infer<typeof StepRunJobPayloadSchema>;
