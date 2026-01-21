import { DbHealthDto } from 'src/prisma/dto/prisma.dto';

export type PrismaServiceMock = {
  $queryRaw: jest.Mock<Promise<unknown>, [TemplateStringsArray, ...unknown[]]>;
  $disconnect: jest.Mock<Promise<void>, []>;
  healthInfo: jest.Mock<Promise<DbHealthDto>, []>;
};

export const createPrismaServiceMock = (): PrismaServiceMock => ({
  $queryRaw: jest.fn<Promise<unknown>, [TemplateStringsArray, ...unknown[]]>(),
  $disconnect: jest.fn<Promise<void>, []>(),
  healthInfo: jest.fn<Promise<DbHealthDto>, []>(),
});
