import { Injectable } from '@nestjs/common';
import type { Prisma, StepRun } from '@prisma/client';
import { PrismaService } from '../prisma.service';

@Injectable()
export class StepRunRepository {
  constructor(private readonly prisma: PrismaService) {}

  create(data: Prisma.StepRunCreateInput): Promise<StepRun> {
    return this.prisma.stepRun.create({ data });
  }

  findManyByWorkflowRun(workflowRunId: string): Promise<StepRun[]> {
    return this.prisma.stepRun.findMany({
      where: { workflowRunId },
      orderBy: { createdAt: 'asc' },
    });
  }

  findById(id: string): Promise<StepRun | null> {
    return this.prisma.stepRun.findUnique({ where: { id } });
  }

  update(id: string, data: Prisma.StepRunUpdateInput): Promise<StepRun> {
    return this.prisma.stepRun.update({ where: { id }, data });
  }
}