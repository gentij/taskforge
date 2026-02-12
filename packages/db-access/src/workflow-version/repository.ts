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

  async findPageByWorkflow(params: {
    workflowId: string;
    page: number;
    pageSize: number;
  }): Promise<{ items: WorkflowVersion[]; total: number }> {
    const skip = (params.page - 1) * params.pageSize;
    const [items, total] = await Promise.all([
      this.prisma.workflowVersion.findMany({
        where: { workflowId: params.workflowId },
        orderBy: { version: 'desc' },
        skip,
        take: params.pageSize,
      }),
      this.prisma.workflowVersion.count({
        where: { workflowId: params.workflowId },
      }),
    ]);

    return { items, total };
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
