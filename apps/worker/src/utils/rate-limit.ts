import type Redis from 'ioredis';

const LUA_FIXED_WINDOW = `
local current = redis.call('INCR', KEYS[1])
if current == 1 then
  redis.call('EXPIRE', KEYS[1], ARGV[2])
end
local ttl = redis.call('TTL', KEYS[1])
return { current, ttl }
`;

export async function checkFixedWindowRateLimit(params: {
  redis: Redis;
  key: string;
  max: number;
  perSeconds: number;
}): Promise<{ allowed: boolean; current: number; ttlSeconds: number }> {
  const { redis, key, max, perSeconds } = params;
  const res = await redis.eval(LUA_FIXED_WINDOW, 1, key, String(max), String(perSeconds));

  const arr = Array.isArray(res) ? res : [];
  const current = typeof arr[0] === 'number' ? arr[0] : Number(arr[0]);
  const ttlSeconds = typeof arr[1] === 'number' ? arr[1] : Number(arr[1]);

  return {
    allowed: Number.isFinite(current) && current <= max,
    current: Number.isFinite(current) ? current : max + 1,
    ttlSeconds: Number.isFinite(ttlSeconds) && ttlSeconds >= 0 ? ttlSeconds : perSeconds,
  };
}
