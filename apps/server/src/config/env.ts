import { z } from 'zod';

const envSchema = z.object({
  NODE_ENV: z
    .enum(['development', 'test', 'production'])
    .default('development'),

  PORT: z
    .string()
    .optional()
    .transform((val) => (val ? Number(val) : 3000))
    .refine((val) => Number.isInteger(val) && val > 0, {
      message: 'PORT must be a positive integer',
    }),

  DATABASE_URL: z
    .string()
    .min(1, 'DATABASE_URL is required')
    .url()
    .refine(
      (url) => url.startsWith('postgresql://') || url.startsWith('postgres://'),
      {
        message: 'DATABASE_URL must be a Postgres connection string',
      },
    ),

  REDIS_URL: z
    .string()
    .min(1, 'REDIS_URL is required')
    .url()
    .refine((url) => url.startsWith('redis://'), {
      message: 'REDIS_URL must be a Redis connection string',
    }),

  TASKFORGE_ADMIN_TOKEN: z
    .string()
    .min(32, 'TASKFORGE_ADMIN_TOKEN must be at least 32 characters'),

  VERSION: z.string(),
});

export function validateEnv(env: NodeJS.ProcessEnv) {
  const parsed = envSchema.safeParse(env);

  if (!parsed.success) {
    console.error('Invalid environment variables:');
    console.error(z.treeifyError(parsed.error));
    throw new Error('Invalid environment variables');
  }

  return parsed.data;
}

export type Env = z.infer<typeof envSchema>;
