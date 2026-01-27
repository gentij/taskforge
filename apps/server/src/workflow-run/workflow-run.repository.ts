import { Injectable } from '@nestjs/common';
import type { Prisma, WorkflowRun } from '@prisma/client';
import { PrismaService } from 'src/prisma/prisma.service';

@Injectable()
export class WorkflowRunRepository {
  constructor(private readonly prisma: PrismaService) {}

  create(data: Prisma.WorkflowRunCreateInput): Promise<WorkflowRun> {
    return this.prisma.workflowRun.create({ data });
  }

  findManyByWorkflow(workflowId: string): Promise<WorkflowRun[]> {
    return this.prisma.workflowRun.findMany({
      where: { workflowId },
      orderBy: { createdAt: 'desc' },
    });
  }

  findById(id: string): Promise<WorkflowRun | null> {
    return this.prisma.workflowRun.findUnique({ where: { id } });
  }

  update(
    id: string,
    data: Prisma.WorkflowRunUpdateInput,
  ): Promise<WorkflowRun> {
    return this.prisma.workflowRun.update({ where: { id }, data });
  }
}
