import { DbHealthDto } from 'src/prisma/dto/prisma.dto';

export type PrismaServiceMock = {
  $queryRaw: jest.Mock<Promise<unknown>, [TemplateStringsArray, ...unknown[]]>;
  $disconnect: jest.Mock<Promise<void>, []>;
  healthInfo: jest.Mock<Promise<DbHealthDto>, []>;
  $transaction: jest.Mock<Promise<unknown>, [(tx: PrismaTxMock) => unknown]>;
};

export const createPrismaServiceMock = (): PrismaServiceMock => ({
  $queryRaw: jest.fn<Promise<unknown>, [TemplateStringsArray, ...unknown[]]>(),
  $disconnect: jest.fn<Promise<void>, []>(),
  healthInfo: jest.fn<Promise<DbHealthDto>, []>(),
  $transaction: jest.fn<Promise<unknown>, [(tx: PrismaTxMock) => unknown]>(),
});

export type PrismaTxMock = {
  workflow: {
    create: jest.Mock;
    update: jest.Mock;
  };
  workflowVersion: {
    create: jest.Mock;
    findFirst: jest.Mock;
  };
};
