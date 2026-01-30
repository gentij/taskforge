import { Injectable } from '@nestjs/common';
import type { Prisma, Trigger } from '@prisma/client';
import { TriggerRepository, WorkflowRepository } from '@taskforge/db-access';
import { AppError } from 'src/common/http/errors/app-error';
import { ErrorDefinitions } from 'src/common/http/errors/error-codes';

@Injectable()
export class TriggerService {
  constructor(
    private readonly repo: TriggerRepository,
    private readonly workflowRepo: WorkflowRepository,
  ) {}

  private async assertWorkflowExists(workflowId: string) {
    const wf = await this.workflowRepo.findById(workflowId);
    if (!wf) throw AppError.notFound(ErrorDefinitions.WORKFLOW.NOT_FOUND);
    return wf;
  }

  async create(params: {
    workflowId: string;
    type: Trigger['type'];
    name?: string;
    isActive?: boolean;
    config?: Prisma.InputJsonValue;
  }): Promise<Trigger> {
    await this.assertWorkflowExists(params.workflowId);

    return this.repo.create({
      workflow: { connect: { id: params.workflowId } },
      type: params.type,
      name: params.name,
      isActive: params.isActive ?? true,
      config: params.config ?? {},
    });
  }

  async list(workflowId: string): Promise<Trigger[]> {
    await this.assertWorkflowExists(workflowId);
    return this.repo.findManyByWorkflow(workflowId);
  }

  async get(workflowId: string, id: string): Promise<Trigger> {
    await this.assertWorkflowExists(workflowId);
    const trigger = await this.repo.findById(id);

    if (!trigger || trigger.workflowId !== workflowId)
      throw AppError.notFound(ErrorDefinitions.TRIGGER.NOT_FOUND);

    return trigger;
  }

  async update(
    workflowId: string,
    id: string,
    patch: {
      name?: string;
      config?: Prisma.InputJsonValue;
      isActive?: boolean;
    },
  ): Promise<Trigger> {
    await this.get(workflowId, id);
    return this.repo.update(id, patch);
  }
}
