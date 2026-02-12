import { Injectable } from '@nestjs/common';
import { LRUCache } from 'lru-cache';
import type { WorkflowRun, WorkflowVersion } from '@prisma/client';

const parsePositiveInt = (value: string | number | undefined, fallback: number) => {
  const parsed = Number(value);
  return Number.isFinite(parsed) && parsed > 0 ? Math.floor(parsed) : fallback;
};

@Injectable()
export class WorkerCacheService {
  private readonly workflowVersionCache = new LRUCache<string, WorkflowVersion>({
    max: parsePositiveInt(
      process.env.WORKER_CACHE_MAX_WORKFLOW_VERSION,
      500,
    ),
    ttl:
      parsePositiveInt(
        process.env.WORKER_CACHE_TTL_WORKFLOW_VERSION_SECONDS,
        300,
      ) * 1000,
  });

  private readonly workflowRunCache = new LRUCache<string, WorkflowRun>({
    max: parsePositiveInt(process.env.WORKER_CACHE_MAX_WORKFLOW_RUN, 1000),
    ttl:
      parsePositiveInt(process.env.WORKER_CACHE_TTL_WORKFLOW_RUN_SECONDS, 60) *
      1000,
  });

  private readonly secretCache = new LRUCache<string, string>({
    max: parsePositiveInt(process.env.WORKER_CACHE_MAX_SECRETS, 500),
    ttl:
      parsePositiveInt(process.env.WORKER_CACHE_TTL_SECRETS_SECONDS, 60) * 1000,
  });

  async getWorkflowVersion(
    id: string,
    loader: () => Promise<WorkflowVersion | null>,
  ): Promise<WorkflowVersion | null> {
    try {
      const cached = this.workflowVersionCache.get(id);
      if (cached) return cached;
    } catch {
      // fail-open
    }

    const value = await loader();
    if (value) {
      try {
        this.workflowVersionCache.set(id, value);
      } catch {
        // fail-open
      }
    }
    return value;
  }

  async getWorkflowRun(
    id: string,
    loader: () => Promise<WorkflowRun | null>,
  ): Promise<WorkflowRun | null> {
    try {
      const cached = this.workflowRunCache.get(id);
      if (cached) return cached;
    } catch {
      // fail-open
    }

    const value = await loader();
    if (value) {
      try {
        this.workflowRunCache.set(id, value);
      } catch {
        // fail-open
      }
    }
    return value;
  }

  getSecret(name: string): string | undefined {
    try {
      return this.secretCache.get(name);
    } catch {
      return undefined;
    }
  }

  setSecret(name: string, value: string): void {
    try {
      this.secretCache.set(name, value);
    } catch {
      // fail-open
    }
  }
}
