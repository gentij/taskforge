import type { Secret } from '@prisma/client';

type SecretOverrides = Partial<Secret>;

export const createSecretFixture = (
  overrides: SecretOverrides = {},
): Secret => {
  const now = new Date('2026-01-23T10:00:00.000Z');

  return {
    id: 'sec_1',
    name: 'API_KEY',
    value: 'secret-value',
    description: 'Primary API key',
    createdAt: now,
    updatedAt: now,
    ...overrides,
  };
};

export const createSecretListFixture = (count = 3): Secret[] =>
  Array.from({ length: count }, (_, i) =>
    createSecretFixture({
      id: `sec_${i + 1}`,
      name: `SECRET_${i + 1}`,
    }),
  );
