import { Injectable } from '@nestjs/common';
import type { Prisma, WorkflowVersion } from '@prisma/client';
import { PrismaService } from '../prisma.service';

@Injectable()
export class WorkflowVersionRepository {
  constructor(private readonly prisma: PrismaService) {}

  create(data: Prisma.WorkflowVersionCreateInput): Promise<WorkflowVersion> {
    return this.prisma.workflowVersion.create({ data });
  }

  findManyByWorkflow(workflowId: string): Promise<WorkflowVersion[]> {
    return this.prisma.workflowVersion.findMany({
      where: { workflowId },
      orderBy: { version: 'desc' },
    });
  }

  findById(id: string): Promise<WorkflowVersion | null> {
    return this.prisma.workflowVersion.findUnique({ where: { id } });
  }

  findByWorkflowAndVersion(
    workflowId: string,
    version: number,
  ): Promise<WorkflowVersion | null> {
    return this.prisma.workflowVersion.findUnique({
      where: {
        workflowId_version: { workflowId, version },
      },
    });
  }

  async getNextVersionNumber(workflowId: string): Promise<number> {
    const latest = await this.prisma.workflowVersion.findFirst({
      where: { workflowId },
      orderBy: { version: 'desc' },
      select: { version: true },
    });

    return (latest?.version ?? 0) + 1;
  }
}