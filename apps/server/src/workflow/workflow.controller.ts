import { Body, Controller, Get, Param, Patch, Post } from '@nestjs/common';
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger';
import { ApiEnvelope } from 'src/common/swagger/envelope/api-envelope.decorator';
import { WorkflowService } from './workflow.service';
import {
  CreateWorkflowReqDto,
  UpdateWorkflowReqDto,
  WorkflowResDto,
  RunWorkflowResDto,
} from './dto/workflow.dto';
import {
  CreateWorkflowVersionReqDto,
  WorkflowVersionResDto,
} from 'src/workflow-version/dto/workflow-version.dto';
import { OrchestrationService } from 'src/core/orchestration.service';

@ApiTags('Workflows')
@ApiBearerAuth('bearer')
@Controller('workflows')
export class WorkflowController {
  constructor(
    private readonly service: WorkflowService,
    private readonly orchestrationService: OrchestrationService,
  ) {}

  @ApiEnvelope(WorkflowResDto, {
    description: 'Create workflow',
    errors: [401, 500],
  })
  @Post()
  create(@Body() body: CreateWorkflowReqDto) {
    return this.service.create(body.name);
  }

  @ApiEnvelope(WorkflowResDto, {
    description: 'List workflows',
    errors: [401, 500],
    isArray: true,
  })
  @Get()
  list() {
    return this.service.list();
  }

  @ApiEnvelope(WorkflowResDto, {
    description: 'Get workflow',
    errors: [401, 404, 500],
  })
  @Get(':id')
  get(@Param('id') id: string) {
    return this.service.get(id);
  }

  @ApiEnvelope(WorkflowResDto, {
    description: 'Update workflow',
    errors: [401, 404, 500],
  })
  @Patch(':id')
  update(@Param('id') id: string, @Body() body: UpdateWorkflowReqDto) {
    return this.service.update(id, body);
  }

  @ApiEnvelope(WorkflowVersionResDto, {
    description: 'Create workflow version',
    errors: [401, 404, 500],
  })
  @Post(':id/versions')
  createVersion(
    @Param('id') id: string,
    @Body() body: CreateWorkflowVersionReqDto,
  ) {
    return this.service.createVersion(id, body.definition);
  }

  @ApiEnvelope(RunWorkflowResDto, {
    description: 'Start workflow manually',
    errors: [401, 404, 500],
  })
  @Post(':id/run')
  async runManual(
    @Param('id') workflowId: string,
    @Body() body?: Record<string, unknown>,
  ) {
    const workflow = await this.service.get(workflowId);

    if (!workflow.latestVersionId) {
      throw new Error('Workflow has no versions');
    }

    const { workflowRunId } = await this.orchestrationService.startWorkflow({
      workflowId,
      workflowVersionId: workflow.latestVersionId,
      eventType: 'MANUAL',
      input: body,
    });

    return { workflowRunId, status: 'QUEUED' };
  }
}
