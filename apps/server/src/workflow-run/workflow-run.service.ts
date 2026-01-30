import { Injectable } from '@nestjs/common';
import type { Prisma, WorkflowRun } from '@prisma/client';
import {
  WorkflowRepository,
  WorkflowRunRepository,
} from '@taskforge/db-access';
import { ErrorDefinitions } from 'src/common/http/errors/error-codes';
import { AppError } from 'src/common/http/errors/app-error';

@Injectable()
export class WorkflowRunService {
  constructor(
    private readonly repo: WorkflowRunRepository,
    private readonly workflowRepo: WorkflowRepository,
  ) {}

  private async assertWorkflowExists(workflowId: string) {
    const wf = await this.workflowRepo.findById(workflowId);
    if (!wf) throw AppError.notFound(ErrorDefinitions.WORKFLOW.NOT_FOUND);
    return wf;
  }

  async create(params: {
    workflowId: string;
    workflowVersionId: string;
    triggerId?: string;
    eventId?: string;
    status?: WorkflowRun['status'];
    input?: Prisma.InputJsonValue;
    output?: Prisma.InputJsonValue;
    startedAt?: Date;
    finishedAt?: Date;
  }): Promise<WorkflowRun> {
    await this.assertWorkflowExists(params.workflowId);

    return this.repo.create({
      workflow: { connect: { id: params.workflowId } },
      workflowVersion: { connect: { id: params.workflowVersionId } },
      trigger: params.triggerId
        ? { connect: { id: params.triggerId } }
        : undefined,
      event: params.eventId ? { connect: { id: params.eventId } } : undefined,
      status: params.status ?? 'QUEUED',
      input: params.input ?? {},
      output: params.output,
      startedAt: params.startedAt,
      finishedAt: params.finishedAt,
    });
  }

  async list(workflowId: string): Promise<WorkflowRun[]> {
    await this.assertWorkflowExists(workflowId);
    return this.repo.findManyByWorkflow(workflowId);
  }

  async get(workflowId: string, id: string): Promise<WorkflowRun> {
    await this.assertWorkflowExists(workflowId);
    const run = await this.repo.findById(id);

    if (!run || run.workflowId !== workflowId)
      throw AppError.notFound(ErrorDefinitions.WORKFLOW_RUN.NOT_FOUND);

    return run;
  }

  async update(
    workflowId: string,
    id: string,
    patch: {
      status?: WorkflowRun['status'];
      input?: Prisma.InputJsonValue;
      output?: Prisma.InputJsonValue;
      startedAt?: Date;
      finishedAt?: Date;
    },
  ): Promise<WorkflowRun> {
    await this.get(workflowId, id);
    return this.repo.update(id, patch);
  }
}
