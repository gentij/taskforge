export type CacheManagerMock = {
  get: jest.Mock<Promise<unknown>, [string]>;
  set: jest.Mock<Promise<void>, [string, unknown]>;
  del: jest.Mock<Promise<void>, [string]>;
};

export const createCacheManagerMock = (): CacheManagerMock => ({
  get: jest.fn<Promise<unknown>, [string]>(),
  set: jest.fn<Promise<void>, [string, unknown]>(),
  del: jest.fn<Promise<void>, [string]>(),
});
