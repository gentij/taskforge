import { createZodDto } from 'nestjs-zod';
import { z } from 'zod';

export const WorkflowRunStatusSchema = z.enum([
  'QUEUED',
  'RUNNING',
  'SUCCEEDED',
  'FAILED',
]);

export const WorkflowRunResSchema = z.object({
  id: z.string(),
  workflowId: z.string(),
  workflowVersionId: z.string(),
  triggerId: z.string().nullable(),
  eventId: z.string().nullable(),
  status: WorkflowRunStatusSchema,
  input: z.unknown(),
  output: z.unknown().nullable(),
  startedAt: z.iso.datetime().nullable(),
  finishedAt: z.iso.datetime().nullable(),
  createdAt: z.iso.datetime(),
  updatedAt: z.iso.datetime(),
});

export class WorkflowRunResDto extends createZodDto(WorkflowRunResSchema) {}
