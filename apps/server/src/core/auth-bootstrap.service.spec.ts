import { Test } from '@nestjs/testing';
import { AuthBootstrapService } from './auth-bootstrap.service';
import { ApiTokenService } from 'src/api-token/api-token.service';
import { CryptoService } from 'src/crypto/crypto.service';
import { PinoLogger } from 'nestjs-pino';

import {
  createApiTokenServiceMock,
  ApiTokenServiceMock,
} from 'test/api-token/api-token.service.mock';
import {
  createCryptoServiceMock,
  CryptoServiceMock,
} from 'test/crypto/crypto.service.mock';
import {
  createPinoLoggerMock,
  PinoLoggerMock,
} from 'test/logger/logger.service.mock';

describe('AuthBootstrapService', () => {
  let service: AuthBootstrapService;
  let apiTokenServiceMock: ApiTokenServiceMock;
  let cryptoServiceMock: CryptoServiceMock;
  let loggerMock: PinoLoggerMock;

  beforeEach(async () => {
    apiTokenServiceMock = createApiTokenServiceMock();
    cryptoServiceMock = createCryptoServiceMock();
    loggerMock = createPinoLoggerMock();

    const moduleRef = await Test.createTestingModule({
      providers: [
        AuthBootstrapService,
        { provide: ApiTokenService, useValue: apiTokenServiceMock },
        { provide: CryptoService, useValue: cryptoServiceMock },
        { provide: PinoLogger, useValue: loggerMock },
      ],
    }).compile();

    service = moduleRef.get(AuthBootstrapService);
  });

  it('skips bootstrap if an active token already exists', async () => {
    apiTokenServiceMock.hasAnyActiveToken.mockResolvedValue(true);

    await service.onModuleInit();

    expect(apiTokenServiceMock.createAdminToken).not.toHaveBeenCalled();
    expect(cryptoServiceMock.generateApiToken).not.toHaveBeenCalled();
    expect(cryptoServiceMock.hashApiToken).not.toHaveBeenCalled();

    expect(loggerMock.info).toHaveBeenCalled();
    expect(loggerMock.warn).not.toHaveBeenCalled();
  });

  it('creates initial admin token if none exists', async () => {
    apiTokenServiceMock.hasAnyActiveToken.mockResolvedValue(false);
    cryptoServiceMock.generateApiToken.mockReturnValue('tf_raw_token_123');
    cryptoServiceMock.hashApiToken.mockReturnValue('hashed_token_abc');

    await service.onModuleInit();

    expect(apiTokenServiceMock.createAdminToken).toHaveBeenCalledWith({
      name: 'initial-admin',
      tokenHash: 'hashed_token_abc',
    });

    const warnCalls = loggerMock.warn.mock.calls.flat().join(' ');
    expect(warnCalls).toContain('tf_raw_token_123');
  });
});
