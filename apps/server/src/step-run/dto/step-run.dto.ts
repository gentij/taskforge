import { createZodDto } from 'nestjs-zod';
import { z } from 'zod';

export const StepRunStatusSchema = z.enum([
  'QUEUED',
  'RUNNING',
  'SUCCEEDED',
  'FAILED',
]);

export const StepRunResSchema = z.object({
  id: z.string(),
  workflowRunId: z.string(),
  stepKey: z.string(),
  status: StepRunStatusSchema,
  attempt: z.number().int(),
  input: z.unknown(),
  output: z.unknown().nullable(),
  error: z.unknown().nullable(),
  logs: z.unknown().nullable(),
  lastErrorAt: z.iso.datetime().nullable(),
  durationMs: z.number().int().nullable(),
  startedAt: z.iso.datetime().nullable(),
  finishedAt: z.iso.datetime().nullable(),
  createdAt: z.iso.datetime(),
  updatedAt: z.iso.datetime(),
});

export class StepRunResDto extends createZodDto(StepRunResSchema) {}
