import { HealthResDto } from 'src/health/dto/health.dto';

export const mockHealthResponse: HealthResDto = {
  status: 'ok',
  version: 'test',
  uptime: 1,
  db: {
    latencyMs: 1,
    ok: true,
  },
};
