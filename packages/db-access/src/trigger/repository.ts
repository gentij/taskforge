import { Injectable } from '@nestjs/common';
import type { Prisma, Trigger } from '@prisma/client';
import { PrismaService } from '../prisma.service';

@Injectable()
export class TriggerRepository {
  constructor(private readonly prisma: PrismaService) {}

  create(data: Prisma.TriggerCreateInput): Promise<Trigger> {
    return this.prisma.trigger.create({ data });
  }

  findManyByWorkflow(workflowId: string): Promise<Trigger[]> {
    return this.prisma.trigger.findMany({
      where: { workflowId },
      orderBy: { createdAt: 'desc' },
    });
  }

  async findPageByWorkflow(params: {
    workflowId: string;
    page: number;
    pageSize: number;
  }): Promise<{ items: Trigger[]; total: number }> {
    const skip = (params.page - 1) * params.pageSize;
    const [items, total] = await Promise.all([
      this.prisma.trigger.findMany({
        where: { workflowId: params.workflowId },
        orderBy: { createdAt: 'desc' },
        skip,
        take: params.pageSize,
      }),
      this.prisma.trigger.count({
        where: { workflowId: params.workflowId },
      }),
    ]);

    return { items, total };
  }

  findById(id: string): Promise<Trigger | null> {
    return this.prisma.trigger.findUnique({ where: { id } });
  }

  update(id: string, data: Prisma.TriggerUpdateInput): Promise<Trigger> {
    return this.prisma.trigger.update({ where: { id }, data });
  }

  softDelete(id: string): Promise<Trigger> {
    return this.prisma.trigger.update({
      where: { id },
      data: { isActive: false },
    });
  }
}
