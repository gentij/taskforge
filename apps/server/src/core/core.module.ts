import { Module } from '@nestjs/common';
import { AppLifecycleService } from './app-lifecycle.service';
import { ApiTokenModule } from 'src/api-token/api-token.module';
import { AuthBootstrapService } from './auth-bootstrap.service';
import { CryptoModule } from 'src/crypto/crypto.module';

@Module({
  imports: [ApiTokenModule, CryptoModule],
  providers: [AppLifecycleService, AuthBootstrapService],
  exports: [ApiTokenModule],
})
export class CoreModule {}
