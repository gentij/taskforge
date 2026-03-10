import type { WorkflowVersion } from '@prisma/client';

export type WorkflowVersionRepositoryMock = {
  findManyByWorkflow: jest.Mock<Promise<WorkflowVersion[]>, [string]>;
  findPageByWorkflow: jest.Mock<
    Promise<{ items: WorkflowVersion[]; total: number }>,
    [
      {
        workflowId: string;
        page: number;
        pageSize: number;
        sortBy: 'version' | 'createdAt';
        sortOrder: 'asc' | 'desc';
      },
    ]
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
      [
        {
          workflowId: string;
          page: number;
          pageSize: number;
          sortBy: 'version' | 'createdAt';
          sortOrder: 'asc' | 'desc';
        },
      ]
    >(),
    findByWorkflowAndVersion: jest.fn<
      Promise<WorkflowVersion | null>,
      [string, number]
    >(),
  });
