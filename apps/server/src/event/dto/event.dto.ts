import { createZodDto } from 'nestjs-zod';
import { z } from 'zod';
import {
  PaginationQuerySchema,
  SortOrderSchema,
} from 'src/common/dto/pagination.dto';

export const EventListSortBySchema = z.enum(['receivedAt', 'createdAt']);
export const EventListQuerySchema = PaginationQuerySchema.extend({
  sortBy: EventListSortBySchema.default('receivedAt'),
  sortOrder: SortOrderSchema.default('desc'),
});
export class EventListQueryDto extends createZodDto(EventListQuerySchema) {}

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
