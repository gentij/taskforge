import type { StepRun, StepRunStatus } from '@prisma/client';

type StepRunOverrides = Partial<StepRun>;

export const createStepRunFixture = (
  overrides: StepRunOverrides = {},
): StepRun => {
  const now = new Date('2026-01-23T10:00:00.000Z');

  return {
    id: 'sr_1',
    workflowRunId: 'wfr_1',
    stepKey: 'step_1',
    status: 'QUEUED' as StepRunStatus,
    attempt: 0,
    input: {},
    output: null,
    error: null,
    logs: null,
    lastErrorAt: null,
    durationMs: null,
    startedAt: null,
    finishedAt: null,
    createdAt: now,
    updatedAt: now,
    ...overrides,
  };
};

export const createStepRunListFixture = (count = 3): StepRun[] =>
  Array.from({ length: count }, (_, i) =>
    createStepRunFixture({
      id: `sr_${i + 1}`,
      stepKey: `step_${i + 1}`,
    }),
  );
