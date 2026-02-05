import { Injectable } from '@nestjs/common';
import type { Prisma, Workflow, WorkflowVersion } from '@prisma/client';
import { WorkflowDefinitionSchema } from '@taskforge/contracts';
import { WorkflowRepository, PrismaService } from '@taskforge/db-access';
import { ErrorDefinitions } from 'src/common/http/errors/error-codes';
import { AppError } from 'src/common/http/errors/app-error';
import {
  getInferredDependencies,
  getExecutionBatchesFromDependencies,
  validateWorkflowDefinitionStrict,
} from './workflow-definition.validator';

@Injectable()
export class WorkflowService {
  constructor(
    private readonly repo: WorkflowRepository,
    private readonly prisma: PrismaService,
  ) {}

  async create(name: string): Promise<Workflow> {
    return this.prisma.$transaction(async (tx) => {
      const workflow = await tx.workflow.create({
        data: { name },
      });

      await tx.trigger.create({
        data: {
          workflowId: workflow.id,
          type: 'MANUAL',
          name: 'Manual',
          isActive: true,
          config: {},
        },
      });

      const v1 = await tx.workflowVersion.create({
        data: {
          workflowId: workflow.id,
          version: 1,
          definition: { steps: [] },
        },
      });

      const updatedWorkflow = await tx.workflow.update({
        where: { id: workflow.id },
        data: { latestVersionId: v1.id },
      });

      return updatedWorkflow;
    });
  }

  async createVersion(
    workflowId: string,
    definition: unknown,
  ): Promise<WorkflowVersion> {
    await this.get(workflowId);

    const normalizedDefinition = WorkflowDefinitionSchema.parse(definition);
    this.validateDefinitionOrThrow(normalizedDefinition);

    return this.prisma.$transaction(async (tx) => {
      const latest = await tx.workflowVersion.findFirst({
        where: { workflowId },
        orderBy: { version: 'desc' },
        select: { version: true },
      });

      const nextVersion = (latest?.version ?? 0) + 1;

      const created = await tx.workflowVersion.create({
        data: {
          workflowId,
          version: nextVersion,
          definition: normalizedDefinition as Prisma.InputJsonValue,
        },
      });

      await tx.workflow.update({
        where: { id: workflowId },
        data: { latestVersionId: created.id },
      });

      return created;
    });
  }

  validateDefinition(definition: unknown) {
    const normalizedDefinition = WorkflowDefinitionSchema.parse(definition);
    const issues = validateWorkflowDefinitionStrict(normalizedDefinition);

    const inferredDependencies = getInferredDependencies(normalizedDefinition);
    let executionBatches: string[][] = [];
    if (issues.length === 0) {
      executionBatches =
        getExecutionBatchesFromDependencies(inferredDependencies);
    }

    return {
      valid: issues.length === 0,
      issues,
      inferredDependencies,
      executionBatches,
    };
  }

  validateDefinitionOrThrow(definition: unknown) {
    const result = this.validateDefinition(definition);
    if (!result.valid) {
      throw AppError.badRequest(
        ErrorDefinitions.COMMON.VALIDATION_ERROR,
        result.issues,
      );
    }

    return result;
  }

  list(): Promise<Workflow[]> {
    return this.repo.findMany();
  }

  async get(id: string): Promise<Workflow> {
    const wf = await this.repo.findById(id);

    if (!wf) throw AppError.notFound(ErrorDefinitions.WORKFLOW.NOT_FOUND);

    return wf;
  }

  async update(
    id: string,
    patch: { name?: string; isActive?: boolean },
  ): Promise<Workflow> {
    await this.get(id);

    return this.repo.update(id, patch);
  }
}
