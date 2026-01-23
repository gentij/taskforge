import { Body, Controller, Get, Param, Patch, Post } from '@nestjs/common';
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger';
import { ApiEnvelope } from 'src/common/swagger/envelope/api-envelope.decorator';
import { WorkflowService } from './workflow.service';
import {
  CreateWorkflowReqDto,
  UpdateWorkflowReqDto,
  WorkflowResDto,
} from './dto/workflow.dto';

@ApiTags('Workflows')
@ApiBearerAuth('bearer')
@Controller('workflows')
export class WorkflowController {
  constructor(private readonly service: WorkflowService) {}

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
}
