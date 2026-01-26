import { createZodDto } from 'nestjs-zod';
import { z } from 'zod';

export const TriggerTypeSchema = z.enum(['MANUAL', 'WEBHOOK', 'CRON']);

export const TriggerResSchema = z.object({
  id: z.string(),
  workflowId: z.string(),
  type: TriggerTypeSchema,
  name: z.string().nullable(),
  isActive: z.boolean(),
  config: z.unknown(),
  createdAt: z.iso.datetime(),
  updatedAt: z.iso.datetime(),
});

export class TriggerResDto extends createZodDto(TriggerResSchema) {}

export const CreateTriggerReqSchema = z.object({
  type: TriggerTypeSchema,
  name: z.string().min(1).max(120).optional(),
  config: z.unknown().default({}),
  isActive: z.boolean().optional(),
});

export class CreateTriggerReqDto extends createZodDto(CreateTriggerReqSchema) {}

export const UpdateTriggerReqSchema = z.object({
  name: z.string().min(1).max(120).optional(),
  config: z.unknown().optional(),
  isActive: z.boolean().optional(),
});

export class UpdateTriggerReqDto extends createZodDto(UpdateTriggerReqSchema) {}
