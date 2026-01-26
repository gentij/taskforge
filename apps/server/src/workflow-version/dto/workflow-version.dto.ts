import { createZodDto } from 'nestjs-zod';
import { z } from 'zod';

export const WorkflowVersionResSchema = z.object({
  id: z.string(),
  workflowId: z.string(),
  version: z.number().int().positive(),
  definition: z.unknown(),
  createdAt: z.string().datetime(),
});
export class WorkflowVersionResDto extends createZodDto(
  WorkflowVersionResSchema,
) {}

export const CreateWorkflowVersionReqSchema = z.object({
  definition: z.unknown().default({ steps: [] }),
});
export class CreateWorkflowVersionReqDto extends createZodDto(
  CreateWorkflowVersionReqSchema,
) {}
