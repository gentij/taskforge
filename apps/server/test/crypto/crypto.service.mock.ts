import { CryptoService } from 'src/crypto/crypto.service';

export type CryptoServiceMock = jest.Mocked<
  Pick<
    CryptoService,
    'generateApiToken' | 'hashApiToken' | 'encryptSecret' | 'decryptSecret'
  >
>;

export const createCryptoServiceMock = (): CryptoServiceMock => ({
  generateApiToken: jest.fn(),
  hashApiToken: jest.fn(),
  encryptSecret: jest.fn((value: string) => value),
  decryptSecret: jest.fn((value: string) => value),
});
