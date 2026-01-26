import { Test } from '@nestjs/testing';
import { WorkflowVersionService } from './workflow-version.service';
import { WorkflowVersionRepository } from './workflow-version.repository';
import { WorkflowRepository } from 'src/workflow/workflow.repository';
import {
  createWorkflowVersionRepositoryMock,
  type WorkflowVersionRepositoryMock,
} from 'test/workflow-version/workflow-version.repository.mock';
import {
  createWorkflowVersionFixture,
  createWorkflowVersionListFixture,
} from 'test/workflow-version/workflow-version.fixtures';
import {
  createWorkflowRepositoryMock,
  type WorkflowRepositoryMock,
} from 'test/workflow/workflow.repository.mock';
import { createWorkflowFixture } from 'test/workflow/workflow.fixtures';
import { AppError } from 'src/common/http/errors/app-error';

describe('WorkflowVersionService', () => {
  let service: WorkflowVersionService;
  let repo: WorkflowVersionRepositoryMock;
  let workflowRepo: WorkflowRepositoryMock;

  beforeEach(async () => {
    repo = createWorkflowVersionRepositoryMock();
    workflowRepo = createWorkflowRepositoryMock();

    const moduleRef = await Test.createTestingModule({
      providers: [
        WorkflowVersionService,
        { provide: WorkflowVersionRepository, useValue: repo },
        { provide: WorkflowRepository, useValue: workflowRepo },
      ],
    }).compile();

    service = moduleRef.get(WorkflowVersionService);
  });

  it('list() returns versions for workflow', async () => {
    const wf = createWorkflowFixture({ id: 'wf_1' });
    const list = createWorkflowVersionListFixture(2);

    workflowRepo.findById.mockResolvedValue(wf);
    repo.findManyByWorkflow.mockResolvedValue(list);

    await expect(service.list('wf_1')).resolves.toBe(list);
    expect(repo.findManyByWorkflow).toHaveBeenCalledWith('wf_1');
  });

  it('list() throws notFound when workflow missing', async () => {
    workflowRepo.findById.mockResolvedValue(null);

    await expect(service.list('missing')).rejects.toBeInstanceOf(AppError);
    expect(repo.findManyByWorkflow).not.toHaveBeenCalled();
  });

  it('get() returns a version when found', async () => {
    const wf = createWorkflowFixture({ id: 'wf_1' });
    const version = createWorkflowVersionFixture({
      workflowId: 'wf_1',
      version: 1,
    });

    workflowRepo.findById.mockResolvedValue(wf);
    repo.findByWorkflowAndVersion.mockResolvedValue(version);

    await expect(service.get('wf_1', 1)).resolves.toBe(version);
    expect(repo.findByWorkflowAndVersion).toHaveBeenCalledWith('wf_1', 1);
  });

  it('get() throws notFound when workflow missing', async () => {
    workflowRepo.findById.mockResolvedValue(null);

    await expect(service.get('missing', 1)).rejects.toBeInstanceOf(AppError);
    expect(repo.findByWorkflowAndVersion).not.toHaveBeenCalled();
  });

  it('get() throws notFound when version missing', async () => {
    const wf = createWorkflowFixture({ id: 'wf_1' });
    workflowRepo.findById.mockResolvedValue(wf);
    repo.findByWorkflowAndVersion.mockResolvedValue(null);

    await expect(service.get('wf_1', 42)).rejects.toBeInstanceOf(AppError);
  });
});
