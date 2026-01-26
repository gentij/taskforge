import { Injectable } from '@nestjs/common';
import type { Prisma, Workflow, WorkflowVersion } from '@prisma/client';
import { WorkflowRepository } from './workflow.repository';
import { ErrorDefinitions } from 'src/common/http/errors/error-codes';
import { PrismaService } from 'src/prisma/prisma.service';
import { AppError } from 'src/common/http/errors/app-error';

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
          definition: definition as Prisma.InputJsonValue,
        },
      });

      await tx.workflow.update({
        where: { id: workflowId },
        data: { latestVersionId: created.id },
      });

      return created;
    });
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
