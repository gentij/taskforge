import { createZodDto } from 'nestjs-zod';
import { z } from 'zod';
import { WorkflowDefinitionSchema } from '@taskforge/contracts';

export const WorkflowVersionResSchema = z.object({
  id: z.string(),
  workflowId: z.string(),
  version: z.number().int().positive(),
  definition: WorkflowDefinitionSchema,
  createdAt: z.string().datetime(),
});
export class WorkflowVersionResDto extends createZodDto(
  WorkflowVersionResSchema,
) {}

export const CreateWorkflowVersionReqSchema = z.object({
  definition: WorkflowDefinitionSchema,
});
export class CreateWorkflowVersionReqDto extends createZodDto(
  CreateWorkflowVersionReqSchema,
) {}
