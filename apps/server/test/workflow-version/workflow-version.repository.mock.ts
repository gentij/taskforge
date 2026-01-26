import type { WorkflowVersion } from '@prisma/client';

export type WorkflowVersionRepositoryMock = {
  findManyByWorkflow: jest.Mock<Promise<WorkflowVersion[]>, [string]>;
  findByWorkflowAndVersion: jest.Mock<
    Promise<WorkflowVersion | null>,
    [string, number]
  >;
};

export const createWorkflowVersionRepositoryMock =
  (): WorkflowVersionRepositoryMock => ({
    findManyByWorkflow: jest.fn<Promise<WorkflowVersion[]>, [string]>(),
    findByWorkflowAndVersion: jest.fn<
      Promise<WorkflowVersion | null>,
      [string, number]
    >(),
  });
