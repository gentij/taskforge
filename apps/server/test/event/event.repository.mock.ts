import type { Event } from '@prisma/client';

export type EventRepositoryMock = {
  create: jest.Mock<Promise<Event>, [any]>;
  findManyByTrigger: jest.Mock<Promise<Event[]>, [string]>;
  findPageByTrigger: jest.Mock<
    Promise<{ items: Event[]; total: number }>,
    [{ triggerId: string; page: number; pageSize: number }]
  >;
  findById: jest.Mock<Promise<Event | null>, [string]>;
};

export const createEventRepositoryMock = (): EventRepositoryMock => ({
  create: jest.fn<Promise<Event>, [any]>(),
  findManyByTrigger: jest.fn<Promise<Event[]>, [string]>(),
  findPageByTrigger: jest.fn<
    Promise<{ items: Event[]; total: number }>,
    [{ triggerId: string; page: number; pageSize: number }]
  >(),
  findById: jest.fn<Promise<Event | null>, [string]>(),
});
