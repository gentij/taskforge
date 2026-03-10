import { Controller, Get, Param, Query } from '@nestjs/common';
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger';
import {
  ApiEnvelope,
  ApiPaginatedEnvelope,
} from 'src/common/swagger/envelope/api-envelope.decorator';
import { WorkflowRunService } from './workflow-run.service';
import {
  WorkflowRunListQueryDto,
  WorkflowRunResDto,
} from './dto/workflow-run.dto';

@ApiTags('Workflow Runs')
@ApiBearerAuth('bearer')
@Controller('workflows/:workflowId/runs')
export class WorkflowRunController {
  constructor(private readonly service: WorkflowRunService) {}

  @ApiPaginatedEnvelope(WorkflowRunResDto, {
    description: 'List workflow runs',
    errors: [401, 404, 500],
  })
  @Get()
  list(
    @Param('workflowId') workflowId: string,
    @Query() query: WorkflowRunListQueryDto,
  ) {
    return this.service.list({ workflowId, ...query });
  }

  @ApiEnvelope(WorkflowRunResDto, {
    description: 'Get workflow run',
    errors: [401, 404, 500],
  })
  @Get(':id')
  get(@Param('workflowId') workflowId: string, @Param('id') id: string) {
    return this.service.get(workflowId, id);
  }
}
