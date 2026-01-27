import type { StepRun } from '@prisma/client';

export type StepRunRepositoryMock = {
  create: jest.Mock<Promise<StepRun>, [any]>;
  findManyByWorkflowRun: jest.Mock<Promise<StepRun[]>, [string]>;
  findById: jest.Mock<Promise<StepRun | null>, [string]>;
  update: jest.Mock<Promise<StepRun>, [string, any]>;
};

export const createStepRunRepositoryMock = (): StepRunRepositoryMock => ({
  create: jest.fn<Promise<StepRun>, [any]>(),
  findManyByWorkflowRun: jest.fn<Promise<StepRun[]>, [string]>(),
  findById: jest.fn<Promise<StepRun | null>, [string]>(),
  update: jest.fn<Promise<StepRun>, [string, any]>(),
});
