import { Controller, Get, Param } from '@nestjs/common';
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger';
import { ApiEnvelope } from 'src/common/swagger/envelope/api-envelope.decorator';
import { StepRunService } from './step-run.service';
import { StepRunResDto } from './dto/step-run.dto';

@ApiTags('Step Runs')
@ApiBearerAuth('bearer')
@Controller('workflows/:workflowId/runs/:runId/steps')
export class StepRunController {
  constructor(private readonly service: StepRunService) {}

  @ApiEnvelope(StepRunResDto, {
    description: 'List step runs',
    isArray: true,
    errors: [401, 404, 500],
  })
  @Get()
  list(@Param('workflowId') workflowId: string, @Param('runId') runId: string) {
    return this.service.list(runId);
  }

  @ApiEnvelope(StepRunResDto, {
    description: 'Get step run',
    errors: [401, 404, 500],
  })
  @Get(':id')
  get(
    @Param('workflowId') workflowId: string,
    @Param('runId') runId: string,
    @Param('id') id: string,
  ) {
    return this.service.get(runId, id);
  }
}
