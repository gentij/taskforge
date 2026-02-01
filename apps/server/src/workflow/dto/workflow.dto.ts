import { createZodDto } from 'nestjs-zod';
import { z } from 'zod';

export const WorkflowResSchema = z.object({
  id: z.string(),
  name: z.string(),
  isActive: z.boolean(),
  createdAt: z.iso.datetime(),
  updatedAt: z.iso.datetime(),
});

export class WorkflowResDto extends createZodDto(WorkflowResSchema) {}

export const CreateWorkflowReqSchema = z.object({
  name: z.string().min(1).max(120),
});
export class CreateWorkflowReqDto extends createZodDto(
  CreateWorkflowReqSchema,
) {}

export const UpdateWorkflowReqSchema = z.object({
  name: z.string().min(1).max(120).optional(),
  isActive: z.boolean().optional(),
});
export class UpdateWorkflowReqDto extends createZodDto(
  UpdateWorkflowReqSchema,
) {}

export const RunWorkflowResSchema = z.object({
  workflowRunId: z.string(),
  status: z.string(),
});
export class RunWorkflowResDto extends createZodDto(RunWorkflowResSchema) {}

export const RunWorkflowOverrideSchema = z
  .object({
    query: z
      .record(z.string(), z.union([z.string(), z.number(), z.boolean()]))
      .optional(),
    body: z.unknown().optional(),
  })
  .strict();

export const RunWorkflowReqSchema = z.object({
  input: z.record(z.string(), z.unknown()).default({}),
  overrides: z.record(z.string(), RunWorkflowOverrideSchema).default({}),
});
export class RunWorkflowReqDto extends createZodDto(RunWorkflowReqSchema) {}
