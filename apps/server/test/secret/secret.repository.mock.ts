import type { Secret } from '@prisma/client';

export type SecretRepositoryMock = {
  create: jest.Mock<Promise<Secret>, [any]>;
  findMany: jest.Mock<Promise<Secret[]>, []>;
  findPage: jest.Mock<
    Promise<{ items: Secret[]; total: number }>,
    [{ page: number; pageSize: number }]
  >;
  findById: jest.Mock<Promise<Secret | null>, [string]>;
  update: jest.Mock<Promise<Secret>, [string, any]>;
  delete: jest.Mock<Promise<Secret>, [string]>;
};

export const createSecretRepositoryMock = (): SecretRepositoryMock => ({
  create: jest.fn<Promise<Secret>, [any]>(),
  findMany: jest.fn<Promise<Secret[]>, []>(),
  findPage: jest.fn<
    Promise<{ items: Secret[]; total: number }>,
    [{ page: number; pageSize: number }]
  >(),
  findById: jest.fn<Promise<Secret | null>, [string]>(),
  update: jest.fn<Promise<Secret>, [string, any]>(),
  delete: jest.fn<Promise<Secret>, [string]>(),
});
