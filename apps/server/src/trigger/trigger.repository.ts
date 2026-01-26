import { Injectable } from '@nestjs/common';
import type { Prisma, Trigger } from '@prisma/client';
import { PrismaService } from 'src/prisma/prisma.service';

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

  findById(id: string): Promise<Trigger | null> {
    return this.prisma.trigger.findUnique({ where: { id } });
  }

  update(id: string, data: Prisma.TriggerUpdateInput): Promise<Trigger> {
    return this.prisma.trigger.update({ where: { id }, data });
  }
}
