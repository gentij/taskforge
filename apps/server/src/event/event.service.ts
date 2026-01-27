import { Injectable } from '@nestjs/common';
import type { Prisma, Event } from '@prisma/client';
import { WorkflowRepository } from 'src/workflow/workflow.repository';
import { TriggerRepository } from 'src/trigger/trigger.repository';
import { ErrorDefinitions } from 'src/common/http/errors/error-codes';
import { AppError } from 'src/common/http/errors/app-error';
import { EventRepository } from './event.repository';

@Injectable()
export class EventService {
  constructor(
    private readonly repo: EventRepository,
    private readonly workflowRepo: WorkflowRepository,
    private readonly triggerRepo: TriggerRepository,
  ) {}

  private async assertTriggerExists(workflowId: string, triggerId: string) {
    const wf = await this.workflowRepo.findById(workflowId);
    if (!wf) throw AppError.notFound(ErrorDefinitions.WORKFLOW.NOT_FOUND);

    const trigger = await this.triggerRepo.findById(triggerId);
    if (!trigger || trigger.workflowId !== workflowId)
      throw AppError.notFound(ErrorDefinitions.TRIGGER.NOT_FOUND);

    return trigger;
  }

  async create(params: {
    triggerId: string;
    type?: string;
    externalId?: string;
    payload?: Prisma.InputJsonValue;
    receivedAt?: Date;
  }): Promise<Event> {
    return this.repo.create({
      trigger: { connect: { id: params.triggerId } },
      type: params.type,
      externalId: params.externalId,
      payload: params.payload ?? {},
      receivedAt: params.receivedAt,
    });
  }

  async list(workflowId: string, triggerId: string): Promise<Event[]> {
    await this.assertTriggerExists(workflowId, triggerId);
    return this.repo.findManyByTrigger(triggerId);
  }

  async get(workflowId: string, triggerId: string, id: string): Promise<Event> {
    await this.assertTriggerExists(workflowId, triggerId);
    const event = await this.repo.findById(id);

    if (!event || event.triggerId !== triggerId)
      throw AppError.notFound(ErrorDefinitions.EVENT.NOT_FOUND);

    return event;
  }
}
