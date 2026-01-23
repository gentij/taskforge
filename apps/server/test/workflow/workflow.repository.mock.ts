import type { Workflow } from '@prisma/client';

export type WorkflowRepositoryMock = {
  create: jest.Mock<Promise<Workflow>, [{ name: string }]>;
  findMany: jest.Mock<Promise<Workflow[]>, []>;
  findById: jest.Mock<Promise<Workflow | null>, [string]>;
  update: jest.Mock<
    Promise<Workflow>,
    [string, { name?: string; isActive?: boolean }]
  >;
};

export const createWorkflowRepositoryMock = (): WorkflowRepositoryMock => ({
  create: jest.fn<Promise<Workflow>, [{ name: string }]>(),
  findMany: jest.fn<Promise<Workflow[]>, []>(),
  findById: jest.fn<Promise<Workflow | null>, [string]>(),
  update: jest.fn<
    Promise<Workflow>,
    [string, { name?: string; isActive?: boolean }]
  >(),
});
