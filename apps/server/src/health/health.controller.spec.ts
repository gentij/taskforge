import { Test, TestingModule } from '@nestjs/testing';
import { HealthController } from './health.controller';
import { HealthService } from './health.service';
import { ConfigService } from '@nestjs/config';
import { mockHealthResponse } from 'test/health/health.mocks';

describe('HealthController', () => {
  let controller: HealthController;
  let service: HealthService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [HealthController],
      providers: [
        {
          provide: HealthService,
          useValue: {
            health: jest.fn().mockReturnValue(mockHealthResponse),
          },
        },
        ConfigService,
      ],
    }).compile();

    controller = module.get<HealthController>(HealthController);
    service = module.get<HealthService>(HealthService);
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });

  it('should call HealthService.health()', async () => {
    const healthSpy = jest.spyOn(service, 'health');

    await controller.health();

    expect(healthSpy).toHaveBeenCalledTimes(1);
  });

  it('should return the health response from the service', () => {
    const result = controller.health();

    expect(result).toEqual(mockHealthResponse);
  });
});
