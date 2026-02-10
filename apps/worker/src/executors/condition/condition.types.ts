import { z } from 'zod';

export const ConditionRequestSpecSchema = z
  .object({
    expr: z.string().min(1),
    assert: z.boolean().optional(),
    message: z.string().min(1).optional(),
    source: z.record(z.string(), z.unknown()).optional(),
  })
  .strict();

export type ConditionRequestSpec = z.infer<typeof ConditionRequestSpecSchema>;

export const ConditionExecutorInputSchema = z.object({
  request: ConditionRequestSpecSchema,
  input: z.unknown().default({}),
});

export type ConditionExecutorInput = z.infer<typeof ConditionExecutorInputSchema>;
