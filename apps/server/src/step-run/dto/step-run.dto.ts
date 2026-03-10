import { createZodDto } from 'nestjs-zod';
import { z } from 'zod';
import {
  PaginationQuerySchema,
  SortOrderSchema,
} from 'src/common/dto/pagination.dto';

export const StepRunListSortBySchema = z.enum(['createdAt', 'updatedAt']);
export const StepRunListQuerySchema = PaginationQuerySchema.extend({
  sortBy: StepRunListSortBySchema.default('createdAt'),
  sortOrder: SortOrderSchema.default('asc'),
});
export class StepRunListQueryDto extends createZodDto(StepRunListQuerySchema) {}

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
  requestOverride: z.unknown().nullable(),
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
