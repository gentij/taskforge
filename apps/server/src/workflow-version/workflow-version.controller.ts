import { Controller, Get, Param, ParseIntPipe } from '@nestjs/common';
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger';
import { ApiEnvelope } from 'src/common/swagger/envelope/api-envelope.decorator';
import { WorkflowVersionService } from './workflow-version.service';
import { WorkflowVersionResDto } from './dto/workflow-version.dto';

@ApiTags('Workflow Versions')
@ApiBearerAuth('bearer')
@Controller('workflows/:workflowId/versions')
export class WorkflowVersionController {
  constructor(private readonly service: WorkflowVersionService) {}

  @ApiEnvelope(WorkflowVersionResDto, {
    description: 'List workflow versions',
    isArray: true,
    errors: [401, 404, 500],
  })
  @Get()
  list(@Param('workflowId') workflowId: string) {
    return this.service.list(workflowId);
  }

  @ApiEnvelope(WorkflowVersionResDto, {
    description: 'Get workflow version',
    errors: [401, 404, 500],
  })
  @Get(':version')
  get(
    @Param('workflowId') workflowId: string,
    @Param('version', ParseIntPipe) version: number,
  ) {
    return this.service.get(workflowId, version);
  }
}
