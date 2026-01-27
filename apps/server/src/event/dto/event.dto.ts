import { createZodDto } from 'nestjs-zod';
import { z } from 'zod';

export const EventResSchema = z.object({
  id: z.string(),
  triggerId: z.string(),
  type: z.string().nullable(),
  externalId: z.string().nullable(),
  payload: z.unknown(),
  receivedAt: z.iso.datetime(),
  createdAt: z.iso.datetime(),
});

export class EventResDto extends createZodDto(EventResSchema) {}
