import { Injectable } from '@nestjs/common';
import type { Prisma, Trigger } from '@prisma/client';
import { TriggerRepository, WorkflowRepository } from '@taskforge/db-access';
import { AppError } from 'src/common/http/errors/app-error';
import { ErrorDefinitions } from 'src/common/http/errors/error-codes';
import { buildPaginationMeta } from 'src/common/pagination/pagination';
import { CronTriggerConfigSchema } from './cron/cron-trigger.types';
import { parseExpression } from 'cron-parser';

const CRON_FIVE_FIELD = /^\s*([^\s]+\s+){4}[^\s]+\s*$/;

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

    if (params.type === 'CRON') {
      const parsed = CronTriggerConfigSchema.safeParse(params.config ?? {});
      if (!parsed.success) {
        throw AppError.badRequest(ErrorDefinitions.COMMON.VALIDATION_ERROR, [
          { field: 'config', message: parsed.error.message },
        ]);
      }

      // Validate cron expression (5-field only) and timezone.
      try {
        if (!CRON_FIVE_FIELD.test(parsed.data.cron)) {
          throw new Error('Cron expression must have 5 fields');
        }
        parseExpression(parsed.data.cron, { tz: parsed.data.timezone });
      } catch (err) {
        const msg = err instanceof Error ? err.message : String(err);
        throw AppError.badRequest(ErrorDefinitions.COMMON.VALIDATION_ERROR, [
          { field: 'config.cron', message: `Invalid cron: ${msg}` },
        ]);
      }
    }

    return this.repo.create({
      workflow: { connect: { id: params.workflowId } },
      type: params.type,
      name: params.name,
      isActive: params.isActive ?? true,
      config: params.config ?? {},
    });
  }

  async list(params: {
    workflowId: string;
    page: number;
    pageSize: number;
  }): Promise<{
    items: Trigger[];
    pagination: ReturnType<typeof buildPaginationMeta>;
  }> {
    await this.assertWorkflowExists(params.workflowId);
    const { items, total } = await this.repo.findPageByWorkflow(params);
    return {
      items,
      pagination: buildPaginationMeta({
        page: params.page,
        pageSize: params.pageSize,
        total,
      }),
    };
  }

  async get(workflowId: string, id: string): Promise<Trigger> {
    await this.assertWorkflowExists(workflowId);
    const trigger = await this.repo.findById(id);

    if (!trigger || trigger.workflowId !== workflowId)
      throw AppError.notFound(ErrorDefinitions.TRIGGER.NOT_FOUND);

    return trigger;
  }

  async delete(workflowId: string, id: string): Promise<Trigger> {
    await this.assertWorkflowExists(workflowId);
    const trigger = await this.repo.findById(id);

    if (!trigger || trigger.workflowId !== workflowId)
      throw AppError.notFound(ErrorDefinitions.TRIGGER.NOT_FOUND);

    return this.repo.softDelete(id);
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

    // If updating config for a CRON trigger, validate it.
    const existing = await this.repo.findById(id);
    if (existing?.type === 'CRON' && patch.config !== undefined) {
      const parsed = CronTriggerConfigSchema.safeParse(patch.config);
      if (!parsed.success) {
        throw AppError.badRequest(ErrorDefinitions.COMMON.VALIDATION_ERROR, [
          { field: 'config', message: parsed.error.message },
        ]);
      }

      try {
        if (!CRON_FIVE_FIELD.test(parsed.data.cron)) {
          throw new Error('Cron expression must have 5 fields');
        }
        parseExpression(parsed.data.cron, { tz: parsed.data.timezone });
      } catch (err) {
        const msg = err instanceof Error ? err.message : String(err);
        throw AppError.badRequest(ErrorDefinitions.COMMON.VALIDATION_ERROR, [
          { field: 'config.cron', message: `Invalid cron: ${msg}` },
        ]);
      }
    }

    return this.repo.update(id, patch);
  }
}
