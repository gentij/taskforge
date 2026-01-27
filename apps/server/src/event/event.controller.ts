import { Controller, Get, Param } from '@nestjs/common';
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger';
import { ApiEnvelope } from 'src/common/swagger/envelope/api-envelope.decorator';
import { EventService } from './event.service';
import { EventResDto } from './dto/event.dto';

@ApiTags('Events')
@ApiBearerAuth('bearer')
@Controller('workflows/:workflowId/triggers/:triggerId/events')
export class EventController {
  constructor(private readonly service: EventService) {}

  @ApiEnvelope(EventResDto, {
    description: 'List events',
    isArray: true,
    errors: [401, 404, 500],
  })
  @Get()
  list(
    @Param('workflowId') workflowId: string,
    @Param('triggerId') triggerId: string,
  ) {
    return this.service.list(workflowId, triggerId);
  }

  @ApiEnvelope(EventResDto, {
    description: 'Get event',
    errors: [401, 404, 500],
  })
  @Get(':id')
  get(
    @Param('workflowId') workflowId: string,
    @Param('triggerId') triggerId: string,
    @Param('id') id: string,
  ) {
    return this.service.get(workflowId, triggerId, id);
  }
}
