import { createZodDto } from 'nestjs-zod';
import { z } from 'zod';
import {
  PaginationQuerySchema,
  SortOrderSchema,
} from 'src/common/dto/pagination.dto';

export const SecretListSortBySchema = z.enum(['createdAt', 'updatedAt']);
export const SecretListQuerySchema = PaginationQuerySchema.extend({
  sortBy: SecretListSortBySchema.default('createdAt'),
  sortOrder: SortOrderSchema.default('desc'),
});
export class SecretListQueryDto extends createZodDto(SecretListQuerySchema) {}

export const SecretResSchema = z.object({
  id: z.string(),
  name: z.string(),
  value: z.string(),
  description: z.string().nullable(),
  createdAt: z.iso.datetime(),
  updatedAt: z.iso.datetime(),
});

export class SecretResDto extends createZodDto(SecretResSchema) {}

export const CreateSecretReqSchema = z.object({
  name: z.string().min(1).max(120),
  value: z.string().min(1).max(5000),
  description: z.string().min(1).max(500).optional(),
});

export class CreateSecretReqDto extends createZodDto(CreateSecretReqSchema) {}

export const UpdateSecretReqSchema = z.object({
  name: z.string().min(1).max(120).optional(),
  value: z.string().min(1).max(5000).optional(),
  description: z.string().min(1).max(500).optional(),
});

export class UpdateSecretReqDto extends createZodDto(UpdateSecretReqSchema) {}
