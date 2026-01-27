import { Module } from '@nestjs/common';
import { AppLifecycleService } from './app-lifecycle.service';
import { ApiTokenModule } from 'src/api-token/api-token.module';
import { AuthBootstrapService } from './auth-bootstrap.service';
import { CryptoModule } from 'src/crypto/crypto.module';
import { AuthModule } from 'src/auth/auth.module';
import { WorkflowModule } from 'src/workflow/workflow.module';
import { WorkflowVersionModule } from 'src/workflow-version/workflow-version.module';
import { TriggerModule } from 'src/trigger/trigger.module';
import { EventModule } from 'src/event/event.module';

@Module({
  imports: [
    ApiTokenModule,
    CryptoModule,
    AuthModule,
    WorkflowModule,
    WorkflowVersionModule,
    TriggerModule,
    EventModule,
  ],
  providers: [AppLifecycleService, AuthBootstrapService],
  exports: [ApiTokenModule],
})
export class CoreModule {}
