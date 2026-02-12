import { Controller, Get, Param, ParseIntPipe, Query } from '@nestjs/common';
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger';
import {
  ApiEnvelope,
  ApiPaginatedEnvelope,
} from 'src/common/swagger/envelope/api-envelope.decorator';
import { WorkflowVersionService } from './workflow-version.service';
import { WorkflowVersionResDto } from './dto/workflow-version.dto';
import { PaginationQueryDto } from 'src/common/dto/pagination.dto';

@ApiTags('Workflow Versions')
@ApiBearerAuth('bearer')
@Controller('workflows/:workflowId/versions')
export class WorkflowVersionController {
  constructor(private readonly service: WorkflowVersionService) {}

  @ApiPaginatedEnvelope(WorkflowVersionResDto, {
    description: 'List workflow versions',
    errors: [401, 404, 500],
  })
  @Get()
  list(
    @Param('workflowId') workflowId: string,
    @Query() query: PaginationQueryDto,
  ) {
    return this.service.list({ workflowId, ...query });
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
