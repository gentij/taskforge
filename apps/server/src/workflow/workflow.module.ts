import { Module } from '@nestjs/common';
import { WorkflowRepository } from './workflow.repository';
import { PrismaModule } from 'src/prisma/prisma.module';
import { WorkflowController } from './workflow.controller';
import { WorkflowService } from './workflow.service';

@Module({
  imports: [PrismaModule],
  controllers: [WorkflowController],
  providers: [WorkflowService, WorkflowRepository],
  exports: [WorkflowService],
})
export class WorkflowModule {}
