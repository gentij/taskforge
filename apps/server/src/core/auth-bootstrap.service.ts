import { Injectable, OnModuleInit } from '@nestjs/common';
import { ApiTokenService } from '../api-token/api-token.service';
import { CryptoService } from 'src/crypto/crypto.service';
import { PinoLogger } from 'nestjs-pino';

@Injectable()
export class AuthBootstrapService implements OnModuleInit {
  constructor(
    private readonly apiTokenService: ApiTokenService,
    private readonly cryptoService: CryptoService,
    private readonly logger: PinoLogger,
  ) {
    this.logger.setContext(AuthBootstrapService.name);
  }

  async onModuleInit() {
    const hasToken = await this.apiTokenService.hasAnyActiveToken();

    if (hasToken) {
      this.logger.info('API token already exists, skipping bootstrap');
      return;
    }

    const rawToken = this.cryptoService.generateApiToken();

    const tokenHash = this.cryptoService.hashApiToken(rawToken);

    await this.apiTokenService.createAdminToken({
      name: 'initial-admin',
      tokenHash,
    });

    this.logger.warn('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    this.logger.warn('No API token found — generated admin token');
    this.logger.warn('');
    this.logger.warn(`API TOKEN: ${rawToken}`);
    this.logger.warn('');
    this.logger.warn('Store this token securely. It will not be shown again.');
    this.logger.warn('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
  }
}
