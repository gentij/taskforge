import { Test, TestingModule } from '@nestjs/testing';
import { HealthService } from './health.service';
import { ConfigService } from '@nestjs/config';
import { HealthResDto } from './dto/health.dto';

describe('HealthService', () => {
  let service: HealthService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [HealthService, ConfigService],
    }).compile();

    service = module.get<HealthService>(HealthService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });

  it('should return health information', () => {
    const result: HealthResDto = service.health();

    expect(result).toEqual(
      expect.objectContaining({
        status: 'ok',
        version: expect.any(String) as string,
        uptime: expect.any(Number) as number,
        timestamp: expect.any(String) as string,
      }),
    );
  });

  it('should return a positive uptime', () => {
    const result: HealthResDto = service.health();

    expect(result.uptime).toBeGreaterThan(0);
  });

  it('should return a valid ISO timestamp', () => {
    const result: HealthResDto = service.health();

    const date = new Date(result.timestamp);
    expect(date.toString()).not.toBe('Invalid Date');
  });
});
