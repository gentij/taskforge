import { Injectable } from '@nestjs/common';
import type { Event, Prisma } from '@prisma/client';
import { PrismaService } from '../prisma.service';

@Injectable()
export class EventRepository {
  constructor(private readonly prisma: PrismaService) {}

  create(data: Prisma.EventCreateInput): Promise<Event> {
    return this.prisma.event.create({ data });
  }

  findManyByTrigger(triggerId: string): Promise<Event[]> {
    return this.prisma.event.findMany({
      where: { triggerId },
      orderBy: { receivedAt: 'desc' },
    });
  }

  async findPageByTrigger(params: {
    triggerId: string;
    page: number;
    pageSize: number;
  }): Promise<{ items: Event[]; total: number }> {
    const skip = (params.page - 1) * params.pageSize;
    const [items, total] = await Promise.all([
      this.prisma.event.findMany({
        where: { triggerId: params.triggerId },
        orderBy: { receivedAt: 'desc' },
        skip,
        take: params.pageSize,
      }),
      this.prisma.event.count({
        where: { triggerId: params.triggerId },
      }),
    ]);

    return { items, total };
  }

  findById(id: string): Promise<Event | null> {
    return this.prisma.event.findUnique({ where: { id } });
  }
}
