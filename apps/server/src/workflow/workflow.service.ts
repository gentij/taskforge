import { Injectable } from '@nestjs/common';
import type { Workflow } from '@prisma/client';
import { WorkflowRepository } from './workflow.repository';
import { AppError } from 'src/common/http/errors/ app-error';
import { ErrorDefinitions } from 'src/common/http/errors/error-codes';

@Injectable()
export class WorkflowService {
  constructor(private readonly repo: WorkflowRepository) {}

  create(name: string): Promise<Workflow> {
    return this.repo.create({ name });
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
