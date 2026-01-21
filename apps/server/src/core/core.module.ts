import { Module } from '@nestjs/common';
import { AppLifecycleService } from './app-lifecycle.service';
import { ApiTokenModule } from 'src/api-token/api-token.module';
import { AuthBootstrapService } from './auth-bootstrap.service';
import { CryptoModule } from 'src/crypto/crypto.module';
import { AuthModule } from 'src/auth/auth.module';

@Module({
  imports: [ApiTokenModule, CryptoModule, AuthModule],
  providers: [AppLifecycleService, AuthBootstrapService],
  exports: [ApiTokenModule],
})
export class CoreModule {}
