import { checkFixedWindowRateLimit } from './rate-limit';

describe('checkFixedWindowRateLimit', () => {
  it('allows when current <= max', async () => {
    const redis = {
      eval: jest.fn().mockResolvedValue([1, 60]),
    };

    const res = await checkFixedWindowRateLimit({
      redis: redis as any,
      key: 'k',
      max: 10,
      perSeconds: 60,
    });

    expect(res.allowed).toBe(true);
    expect(res.current).toBe(1);
  });

  it('denies when current > max', async () => {
    const redis = {
      eval: jest.fn().mockResolvedValue([11, 10]),
    };

    const res = await checkFixedWindowRateLimit({
      redis: redis as any,
      key: 'k',
      max: 10,
      perSeconds: 60,
    });

    expect(res.allowed).toBe(false);
    expect(res.current).toBe(11);
    expect(res.ttlSeconds).toBe(10);
  });
});
