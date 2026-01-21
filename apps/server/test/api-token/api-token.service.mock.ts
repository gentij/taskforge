import { ApiToken } from '@prisma/client';
import { ApiTokenService } from 'src/api-token/api-token.service';

export type ApiTokenServiceMock = jest.Mocked<
  Pick<
    ApiTokenService,
    'hasAnyActiveToken' | 'createAdminToken' | 'validateTokenHash'
  >
>;

export const createApiTokenServiceMock = (): ApiTokenServiceMock => ({
  hasAnyActiveToken: jest.fn(),
  createAdminToken: jest.fn(),
  validateTokenHash: jest.fn(),
});

type ApiTokenOverrides = Partial<ApiToken>;

export const createApiTokenFixture = (
  overrides: ApiTokenOverrides = {},
): ApiToken => ({
  id: 'token-id',
  name: 'api-token',
  tokenHash: 'hashed-token',
  scopes: [],
  createdAt: new Date('2025-01-01'),
  lastUsedAt: null,
  revokedAt: null,
  ...overrides,
});

export const createRevokedApiTokenFixture = (
  overrides: ApiTokenOverrides = {},
): ApiToken =>
  createApiTokenFixture({
    revokedAt: new Date('2025-01-02'),
    ...overrides,
  });

export const createApiTokenListFixture = (
  count = 1,
  overrides: ApiTokenOverrides = {},
): ApiToken[] =>
  Array.from({ length: count }, (_, i) =>
    createApiTokenFixture({
      id: `token-${i + 1}`,
      ...overrides,
    }),
  );
