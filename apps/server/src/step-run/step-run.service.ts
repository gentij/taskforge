import { Injectable } from '@nestjs/common';
import type { Prisma, StepRun } from '@prisma/client';
import { WorkflowRunRepository } from 'src/workflow-run/workflow-run.repository';
import { ErrorDefinitions } from 'src/common/http/errors/error-codes';
import { AppError } from 'src/common/http/errors/app-error';
import { StepRunRepository } from './step-run.repository';

@Injectable()
export class StepRunService {
  constructor(
    private readonly repo: StepRunRepository,
    private readonly runRepo: WorkflowRunRepository,
  ) {}

  private async assertWorkflowRunExists(workflowRunId: string) {
    const run = await this.runRepo.findById(workflowRunId);
    if (!run) throw AppError.notFound(ErrorDefinitions.WORKFLOW_RUN.NOT_FOUND);
    return run;
  }

  async create(params: {
    workflowRunId: string;
    stepKey: string;
    status?: StepRun['status'];
    attempt?: number;
    input?: Prisma.InputJsonValue;
    output?: Prisma.InputJsonValue;
    error?: Prisma.InputJsonValue;
    logs?: Prisma.InputJsonValue;
    lastErrorAt?: Date;
    durationMs?: number;
    startedAt?: Date;
    finishedAt?: Date;
  }): Promise<StepRun> {
    await this.assertWorkflowRunExists(params.workflowRunId);

    return this.repo.create({
      workflowRun: { connect: { id: params.workflowRunId } },
      stepKey: params.stepKey,
      status: params.status ?? 'QUEUED',
      attempt: params.attempt ?? 0,
      input: params.input ?? {},
      output: params.output,
      error: params.error,
      logs: params.logs,
      lastErrorAt: params.lastErrorAt,
      durationMs: params.durationMs,
      startedAt: params.startedAt,
      finishedAt: params.finishedAt,
    });
  }

  async list(workflowRunId: string): Promise<StepRun[]> {
    await this.assertWorkflowRunExists(workflowRunId);
    return this.repo.findManyByWorkflowRun(workflowRunId);
  }

  async get(workflowRunId: string, id: string): Promise<StepRun> {
    await this.assertWorkflowRunExists(workflowRunId);
    const step = await this.repo.findById(id);

    if (!step || step.workflowRunId !== workflowRunId)
      throw AppError.notFound(ErrorDefinitions.STEP_RUN.NOT_FOUND);

    return step;
  }

  async update(
    workflowRunId: string,
    id: string,
    patch: {
      status?: StepRun['status'];
      attempt?: number;
      input?: Prisma.InputJsonValue;
      output?: Prisma.InputJsonValue;
      error?: Prisma.InputJsonValue;
      logs?: Prisma.InputJsonValue;
      lastErrorAt?: Date;
      durationMs?: number;
      startedAt?: Date;
      finishedAt?: Date;
    },
  ): Promise<StepRun> {
    await this.get(workflowRunId, id);
    return this.repo.update(id, patch);
  }
}
