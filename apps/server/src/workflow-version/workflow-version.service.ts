import { Inject, Injectable } from '@nestjs/common';
import {
  WorkflowVersionRepository,
  WorkflowRepository,
} from '@taskforge/db-access';
import { ErrorDefinitions } from 'src/common/http/errors/error-codes';
import { AppError } from 'src/common/http/errors/app-error';
import { CACHE_MANAGER } from '@nestjs/cache-manager';
import type { Cache } from 'cache-manager';
import { cacheKeys } from 'src/cache/cache-keys';

@Injectable()
export class WorkflowVersionService {
  constructor(
    private readonly workflowRepo: WorkflowRepository,
    private readonly repo: WorkflowVersionRepository,
    @Inject(CACHE_MANAGER) private readonly cache: Cache,
  ) {}

  private async assertWorkflowExists(workflowId: string) {
    const wf = await this.workflowRepo.findById(workflowId);
    if (!wf) throw AppError.notFound(ErrorDefinitions.WORKFLOW.NOT_FOUND);
    return wf;
  }

  async list(workflowId: string) {
    await this.assertWorkflowExists(workflowId);

    const key = cacheKeys.workflowVersionList(workflowId);
    try {
      const cached = await this.cache.get(key);
      if (cached) return cached;
    } catch {
      // fail-open: cache errors should not break API
    }

    const versions = await this.repo.findManyByWorkflow(workflowId);
    try {
      await this.cache.set(key, versions);
    } catch {
      // fail-open: cache errors should not break API
    }
    return versions;
  }

  async get(workflowId: string, version: number) {
    await this.assertWorkflowExists(workflowId);

    const key = cacheKeys.workflowVersionGet(workflowId, version);
    try {
      const cached = await this.cache.get(key);
      if (cached) return cached;
    } catch {
      // fail-open: cache errors should not break API
    }

    const v = await this.repo.findByWorkflowAndVersion(workflowId, version);
    if (!v)
      throw AppError.notFound(ErrorDefinitions.WORKFLOW.VERSION_NOT_FOUND);

    try {
      await this.cache.set(key, v);
    } catch {
      // fail-open: cache errors should not break API
    }
    return v;
  }
}
