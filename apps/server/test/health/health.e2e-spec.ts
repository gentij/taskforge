import { Test, TestingModule } from '@nestjs/testing';
import { AppModule } from '../../src/app.module';

import {
  FastifyAdapter,
  NestFastifyApplication,
} from '@nestjs/platform-fastify';
import request from 'supertest';

describe('Health (e2e)', () => {
  let app: NestFastifyApplication;

  beforeAll(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [AppModule],
    }).compile();

    app = moduleFixture.createNestApplication<NestFastifyApplication>(
      new FastifyAdapter(),
    );

    await app.init();
    // Required for Fastify in some setups to ensure the underlying server is ready
    await app.getHttpAdapter().getInstance().ready();
  });

  afterAll(async () => {
    await app.close();
  });

  it('GET /health should return ok + metadata', async () => {
    const res = await request(app.getHttpServer()).get('/health').expect(200);

    expect(res.body).toEqual(
      expect.objectContaining({
        status: 'ok',
        version: expect.any(String) as string,
        uptime: expect.any(Number) as number,
        db: {
          latencyMs: expect.any(Number) as number,
          ok: true,
        },
      }),
    );
  });
});
