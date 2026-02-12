import type { WorkflowRun } from '@prisma/client';

export type WorkflowRunRepositoryMock = {
  create: jest.Mock<Promise<WorkflowRun>, [any]>;
  findManyByWorkflow: jest.Mock<Promise<WorkflowRun[]>, [string]>;
  findPageByWorkflow: jest.Mock<
    Promise<{ items: WorkflowRun[]; total: number }>,
    [{ workflowId: string; page: number; pageSize: number }]
  >;
  findById: jest.Mock<Promise<WorkflowRun | null>, [string]>;
  update: jest.Mock<Promise<WorkflowRun>, [string, any]>;
};

export const createWorkflowRunRepositoryMock =
  (): WorkflowRunRepositoryMock => ({
    create: jest.fn<Promise<WorkflowRun>, [any]>(),
    findManyByWorkflow: jest.fn<Promise<WorkflowRun[]>, [string]>(),
    findPageByWorkflow: jest.fn<
      Promise<{ items: WorkflowRun[]; total: number }>,
      [{ workflowId: string; page: number; pageSize: number }]
    >(),
    findById: jest.fn<Promise<WorkflowRun | null>, [string]>(),
    update: jest.fn<Promise<WorkflowRun>, [string, any]>(),
  });
