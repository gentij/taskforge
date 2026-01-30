import { Injectable } from '@nestjs/common';
import {
  WorkflowVersionRepository,
  WorkflowRepository,
} from '@taskforge/db-access';
import { ErrorDefinitions } from 'src/common/http/errors/error-codes';
import { AppError } from 'src/common/http/errors/app-error';

@Injectable()
export class WorkflowVersionService {
  constructor(
    private readonly workflowRepo: WorkflowRepository,
    private readonly repo: WorkflowVersionRepository,
  ) {}

  private async assertWorkflowExists(workflowId: string) {
    const wf = await this.workflowRepo.findById(workflowId);
    if (!wf) throw AppError.notFound(ErrorDefinitions.WORKFLOW.NOT_FOUND);
    return wf;
  }

  async list(workflowId: string) {
    await this.assertWorkflowExists(workflowId);
    return this.repo.findManyByWorkflow(workflowId);
  }

  async get(workflowId: string, version: number) {
    await this.assertWorkflowExists(workflowId);

    const v = await this.repo.findByWorkflowAndVersion(workflowId, version);
    if (!v)
      throw AppError.notFound(ErrorDefinitions.WORKFLOW.VERSION_NOT_FOUND);
    return v;
  }
}
