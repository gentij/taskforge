import { Module } from '@nestjs/common';
import { PrismaModule } from 'src/prisma/prisma.module';
import { WorkflowModule } from 'src/workflow/workflow.module';
import { TriggerController } from './trigger.controller';
import { TriggerRepository } from '@taskforge/db-access';
import { TriggerService } from './trigger.service';

@Module({
  imports: [PrismaModule, WorkflowModule],
  controllers: [TriggerController],
  providers: [TriggerService, TriggerRepository],
  exports: [TriggerService, TriggerRepository],
})
export class TriggerModule {}
