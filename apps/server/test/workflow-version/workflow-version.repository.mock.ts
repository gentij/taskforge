import type { WorkflowVersion } from '@prisma/client';

export type WorkflowVersionRepositoryMock = {
  findManyByWorkflow: jest.Mock<Promise<WorkflowVersion[]>, [string]>;
  findPageByWorkflow: jest.Mock<
    Promise<{ items: WorkflowVersion[]; total: number }>,
    [{ workflowId: string; page: number; pageSize: number }]
  >;
  findByWorkflowAndVersion: jest.Mock<
    Promise<WorkflowVersion | null>,
    [string, number]
  >;
};

export const createWorkflowVersionRepositoryMock =
  (): WorkflowVersionRepositoryMock => ({
    findManyByWorkflow: jest.fn<Promise<WorkflowVersion[]>, [string]>(),
    findPageByWorkflow: jest.fn<
      Promise<{ items: WorkflowVersion[]; total: number }>,
      [{ workflowId: string; page: number; pageSize: number }]
    >(),
    findByWorkflowAndVersion: jest.fn<
      Promise<WorkflowVersion | null>,
      [string, number]
    >(),
  });
