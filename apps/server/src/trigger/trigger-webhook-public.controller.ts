import { Body, Controller, Param, Post } from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { AppError } from 'src/common/http/errors/app-error';
import { ErrorDefinitions } from 'src/common/http/errors/error-codes';
import { ApiEnvelope } from 'src/common/swagger/envelope/api-envelope.decorator';
import { OrchestrationService } from 'src/core/orchestration.service';
import { WorkflowService } from 'src/workflow/workflow.service';
import { Public } from 'src/auth/public.decorator';
import { TriggerService } from './trigger.service';
import { TriggerWebhookIngressResDto } from './dto/trigger.dto';

@ApiTags('Webhooks')
@Controller('hooks')
export class TriggerWebhookPublicController {
  constructor(
    private readonly triggerService: TriggerService,
    private readonly orchestrationService: OrchestrationService,
    private readonly workflowService: WorkflowService,
  ) {}

  @ApiEnvelope(TriggerWebhookIngressResDto, {
    description: 'Public webhook ingress',
    errors: [400, 401, 404, 500],
  })
  @Public()
  @Post(':workflowId/:triggerId/:webhookKey')
  async handleWebhook(
    @Param('workflowId') workflowId: string,
    @Param('triggerId') triggerId: string,
    @Param('webhookKey') webhookKey: string,
    @Body() body: unknown,
  ): Promise<{ status: 'accepted' | 'trigger_inactive' }> {
    const trigger = await this.triggerService.get(workflowId, triggerId);
    this.triggerService.assertWebhookTriggerType(trigger);

    if (!this.triggerService.hasWebhookKey(trigger)) {
      throw AppError.badRequest(
        ErrorDefinitions.TRIGGER.WEBHOOK_KEY_NOT_CONFIGURED,
      );
    }

    if (!this.triggerService.hasValidWebhookKey(trigger, webhookKey)) {
      throw AppError.unauthorized([
        { field: 'webhookKey', message: 'Invalid webhook key' },
      ]);
    }

    if (!trigger.isActive) {
      return { status: 'trigger_inactive' };
    }

    const workflow = await this.workflowService.get(workflowId);
    if (!workflow.latestVersionId) {
      throw new Error('Workflow has no versions');
    }

    const input = normalizeWebhookInput(body);

    await this.orchestrationService.startWorkflow({
      workflowId,
      workflowVersionId: workflow.latestVersionId,
      triggerId,
      eventType: 'WEBHOOK',
      eventPayload: input,
      input,
      overrides: {},
    });

    return { status: 'accepted' };
  }
}

function normalizeWebhookInput(payload: unknown): Record<string, unknown> {
  if (!payload || typeof payload !== 'object' || Array.isArray(payload)) {
    return payload === undefined ? {} : { payload };
  }

  return payload as Record<string, unknown>;
}
