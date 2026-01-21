import { CryptoService } from 'src/crypto/crypto.service';

export type CryptoServiceMock = jest.Mocked<
  Pick<CryptoService, 'generateApiToken' | 'hashApiToken'>
>;

export const createCryptoServiceMock = (): CryptoServiceMock => ({
  generateApiToken: jest.fn(),
  hashApiToken: jest.fn(),
});
