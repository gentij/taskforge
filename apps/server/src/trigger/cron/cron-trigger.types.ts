import { z } from 'zod';

export const CronTriggerConfigSchema = z
  .object({
    cron: z.string().min(1),
    timezone: z.string().min(1).optional().default('UTC'),
    input: z.record(z.string(), z.unknown()).optional().default({}),
  })
  .strict();

export type CronTriggerConfig = z.infer<typeof CronTriggerConfigSchema>;
