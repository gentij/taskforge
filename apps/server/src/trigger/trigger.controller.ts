import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Patch,
  Post,
  Query,
} from '@nestjs/common';
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger';
import {
  ApiEnvelope,
  ApiPaginatedEnvelope,
} from 'src/common/swagger/envelope/api-envelope.decorator';
import type { Prisma } from '@prisma/client';
import { TriggerService } from './trigger.service';
import {
  CreateTriggerReqDto,
  TriggerResDto,
  UpdateTriggerReqDto,
} from './dto/trigger.dto';
import { OrchestrationService } from 'src/core/orchestration.service';
import { WorkflowService } from 'src/workflow/workflow.service';
import { RunWorkflowReqDto } from 'src/workflow/dto/workflow.dto';
import { PaginationQueryDto } from 'src/common/dto/pagination.dto';

@ApiTags('Triggers')
@ApiBearerAuth('bearer')
@Controller('workflows/:workflowId/triggers')
export class TriggerController {
  constructor(
    private readonly service: TriggerService,
    private readonly orchestrationService: OrchestrationService,
    private readonly workflowService: WorkflowService,
  ) {}

  @ApiEnvelope(TriggerResDto, {
    description: 'Create trigger',
    errors: [401, 404, 500],
  })
  @Post()
  create(
    @Param('workflowId') workflowId: string,
    @Body() body: CreateTriggerReqDto,
  ) {
    return this.service.create({
      workflowId,
      type: body.type,
      name: body.name,
      isActive: body.isActive,
      config: body.config as Prisma.InputJsonValue,
    });
  }

  @ApiPaginatedEnvelope(TriggerResDto, {
    description: 'List triggers',
    errors: [401, 404, 500],
  })
  @Get()
  list(
    @Param('workflowId') workflowId: string,
    @Query() query: PaginationQueryDto,
  ) {
    return this.service.list({ workflowId, ...query });
  }

  @ApiEnvelope(TriggerResDto, {
    description: 'Get trigger',
    errors: [401, 404, 500],
  })
  @Get(':id')
  get(@Param('workflowId') workflowId: string, @Param('id') id: string) {
    return this.service.get(workflowId, id);
  }

  @ApiEnvelope(TriggerResDto, {
    description: 'Update trigger',
    errors: [401, 404, 500],
  })
  @Patch(':id')
  update(
    @Param('workflowId') workflowId: string,
    @Param('id') id: string,
    @Body() body: UpdateTriggerReqDto,
  ) {
    return this.service.update(workflowId, id, {
      name: body.name,
      isActive: body.isActive,
      config:
        body.config === undefined
          ? undefined
          : (body.config as Prisma.InputJsonValue),
    });
  }

  @ApiEnvelope(TriggerResDto, {
    description: 'Delete trigger (soft)',
    errors: [401, 404, 500],
  })
  @Delete(':id')
  delete(@Param('workflowId') workflowId: string, @Param('id') id: string) {
    return this.service.delete(workflowId, id);
  }

  @ApiTags('Triggers')
  @Post(':id/webhook')
  async handleWebhook(
    @Param('workflowId') workflowId: string,
    @Param('id') triggerId: string,
    @Body() body: RunWorkflowReqDto,
  ) {
    const input = body.input ?? {};
    const overrides = body.overrides ?? {};

    const trigger = await this.service.get(workflowId, triggerId);

    if (!trigger.isActive) {
      return { status: 'trigger_inactive' };
    }

    const workflow = await this.workflowService.get(workflowId);

    if (!workflow.latestVersionId) {
      throw new Error('Workflow has no versions');
    }

    await this.orchestrationService.startWorkflow({
      workflowId,
      workflowVersionId: workflow.latestVersionId,
      triggerId,
      eventType: 'WEBHOOK',
      eventPayload: input,
      input,
      overrides,
    });

    return { status: 'accepted' };
  }
}
