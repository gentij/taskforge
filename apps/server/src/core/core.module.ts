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
import { WorkflowRunModule } from 'src/workflow-run/workflow-run.module';
import { StepRunModule } from 'src/step-run/step-run.module';
import { SecretModule } from 'src/secret/secret.module';
import { QueueModule } from 'src/queue/queue.module';
import { OrchestrationModule } from './orchestration.module';

@Module({
  imports: [
    ApiTokenModule,
    CryptoModule,
    AuthModule,
    WorkflowModule,
    WorkflowVersionModule,
    TriggerModule,
    EventModule,
    WorkflowRunModule,
    StepRunModule,
    SecretModule,
    QueueModule,
    OrchestrationModule,
  ],
  providers: [AppLifecycleService, AuthBootstrapService],
  exports: [ApiTokenModule, OrchestrationModule],
})
export class CoreModule {}
