import { z } from 'zod';

export const TransformRequestSpecSchema = z
  .object({
    source: z.record(z.string(), z.unknown()).optional(),
    output: z.unknown(),
  })
  .strict();

export type TransformRequestSpec = z.infer<typeof TransformRequestSpecSchema>;

export const TransformExecutorInputSchema = z.object({
  request: TransformRequestSpecSchema,
  input: z.unknown().default({}),
});

export type TransformExecutorInput = z.infer<typeof TransformExecutorInputSchema>;
