export type PrismaStepRunModelMock = {
  update: jest.Mock;
  findUnique: jest.Mock;
  findMany: jest.Mock;
};

export type PrismaWorkflowRunModelMock = {
  update: jest.Mock;
  updateMany: jest.Mock;
};

export type PrismaWorkflowVersionModelMock = {
  findUnique: jest.Mock;
};

export type PrismaServiceMock = {
  stepRun: PrismaStepRunModelMock;
  workflowRun: PrismaWorkflowRunModelMock;
  workflowVersion: PrismaWorkflowVersionModelMock;
};

export const createPrismaServiceMock = (): PrismaServiceMock => ({
  stepRun: {
    update: jest.fn(),
    findUnique: jest.fn(),
    findMany: jest.fn(),
  },
  workflowRun: {
    update: jest.fn(),
    updateMany: jest.fn(),
  },
  workflowVersion: {
    findUnique: jest.fn(),
  },
});
