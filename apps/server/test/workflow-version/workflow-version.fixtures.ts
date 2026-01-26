import type { WorkflowVersion } from '@prisma/client';

type WorkflowVersionOverrides = Partial<WorkflowVersion>;

export const createWorkflowVersionFixture = (
  overrides: WorkflowVersionOverrides = {},
): WorkflowVersion => {
  const now = new Date('2026-01-23T10:00:00.000Z');

  return {
    id: 'wfv_1',
    workflowId: 'wf_1',
    version: 1,
    definition: { steps: [] },
    createdAt: now,
    ...overrides,
  };
};

export const createWorkflowVersionListFixture = (
  count = 3,
): WorkflowVersion[] =>
  Array.from({ length: count }, (_, i) =>
    createWorkflowVersionFixture({
      id: `wfv_${i + 1}`,
      version: i + 1,
    }),
  );
