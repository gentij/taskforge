import type { Trigger } from '@prisma/client';

export type TriggerRepositoryMock = {
  create: jest.Mock<Promise<Trigger>, [any]>;
  findManyByWorkflow: jest.Mock<Promise<Trigger[]>, [string]>;
  findById: jest.Mock<Promise<Trigger | null>, [string]>;
  update: jest.Mock<Promise<Trigger>, [string, any]>;
};

export const createTriggerRepositoryMock = (): TriggerRepositoryMock => ({
  create: jest.fn<Promise<Trigger>, [any]>(),
  findManyByWorkflow: jest.fn<Promise<Trigger[]>, [string]>(),
  findById: jest.fn<Promise<Trigger | null>, [string]>(),
  update: jest.fn<Promise<Trigger>, [string, any]>(),
});
