import { PrismaService } from './prisma.service';
import { ConfigService } from '@nestjs/config';
import {
  createPrismaServiceMock,
  PrismaServiceMock,
} from 'test/prisma/prisma.mocks';

type PrismaClientLike = {
  $queryRaw: PrismaServiceMock['$queryRaw'];
  $disconnect: PrismaServiceMock['$disconnect'];
};

describe('PrismaService', () => {
  let service: PrismaService;
  let prismaMock: PrismaServiceMock;

  const configServiceMock = { get: jest.fn() } as unknown as ConfigService;

  beforeEach(() => {
    jest.resetAllMocks();

    prismaMock = createPrismaServiceMock();
    service = new PrismaService(configServiceMock);

    const client = service as unknown as PrismaClientLike;
    client.$queryRaw = prismaMock.$queryRaw;
    client.$disconnect = prismaMock.$disconnect;
  });

  it('healthInfo returns ok=true when SELECT 1 succeeds', async () => {
    prismaMock.$queryRaw.mockResolvedValueOnce(1);

    const result = await service.healthInfo();

    expect(prismaMock.$queryRaw).toHaveBeenCalledTimes(1);
    expect(result.ok).toBe(true);
    expect(result.latencyMs).toEqual(expect.any(Number));
  });

  it('healthInfo returns ok=false when SELECT 1 fails', async () => {
    prismaMock.$queryRaw.mockRejectedValueOnce(new Error('ECONNREFUSED'));

    const result = await service.healthInfo();

    expect(result.ok).toBe(false);
    expect(result.error).toContain('ECONNREFUSED');
  });

  it('onModuleDestroy disconnects', async () => {
    prismaMock.$disconnect.mockResolvedValueOnce(undefined);

    await service.onModuleDestroy();

    expect(prismaMock.$disconnect).toHaveBeenCalledTimes(1);
  });
});
