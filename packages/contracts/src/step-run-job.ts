import { z } from 'zod';

export const StepRunJobPayloadSchema = z.object({
  workflowRunId: z.string().min(1),
  stepRunId: z.string().min(1),
  stepKey: z.string().min(1),
  workflowVersionId: z.string().min(1),
  input: z.unknown().default({}),
});

export type StepRunJobPayload = z.infer<typeof StepRunJobPayloadSchema>;
