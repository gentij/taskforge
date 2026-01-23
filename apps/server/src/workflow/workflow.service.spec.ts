import { Test } from '@nestjs/testing';
import { WorkflowService } from './workflow.service';
import { WorkflowRepository } from './workflow.repository';
import {
  createWorkflowRepositoryMock,
  type WorkflowRepositoryMock,
} from 'test/workflow/workflow.repository.mock';
import {
  createWorkflowFixture,
  createWorkflowListFixture,
} from 'test/workflow/workflow.fixtures';
import { AppError } from 'src/common/http/errors/ app-error';

describe('WorkflowService', () => {
  let service: WorkflowService;
  let repo: WorkflowRepositoryMock;

  beforeEach(async () => {
    repo = createWorkflowRepositoryMock();

    const moduleRef = await Test.createTestingModule({
      providers: [
        WorkflowService,
        { provide: WorkflowRepository, useValue: repo },
      ],
    }).compile();

    service = moduleRef.get(WorkflowService);
  });

  it('create() creates a workflow', async () => {
    const wf = createWorkflowFixture({ name: 'My WF' });
    repo.create.mockResolvedValue(wf);

    await expect(service.create('My WF')).resolves.toBe(wf);
    expect(repo.create).toHaveBeenCalledWith({ name: 'My WF' });
  });

  it('list() returns workflows', async () => {
    const list = createWorkflowListFixture(2);
    repo.findMany.mockResolvedValue(list);

    await expect(service.list()).resolves.toBe(list);
    expect(repo.findMany).toHaveBeenCalledTimes(1);
  });

  it('get() returns workflow when found', async () => {
    const wf = createWorkflowFixture({ id: 'wf_x' });
    repo.findById.mockResolvedValue(wf);

    await expect(service.get('wf_x')).resolves.toBe(wf);
    expect(repo.findById).toHaveBeenCalledWith('wf_x');
  });

  it('get() throws AppError.notFound when missing', async () => {
    repo.findById.mockResolvedValue(null);

    await expect(service.get('missing')).rejects.toBeInstanceOf(AppError);
  });

  it('update() updates after existence check', async () => {
    const existing = createWorkflowFixture({ id: 'wf_1', name: 'Old' });
    const updated = createWorkflowFixture({ id: 'wf_1', name: 'New' });

    repo.findById.mockResolvedValue(existing);
    repo.update.mockResolvedValue(updated);

    await expect(service.update('wf_1', { name: 'New' })).resolves.toBe(
      updated,
    );

    expect(repo.findById).toHaveBeenCalledWith('wf_1');
    expect(repo.update).toHaveBeenCalledWith('wf_1', { name: 'New' });
  });

  it('update() throws notFound when workflow missing', async () => {
    repo.findById.mockResolvedValue(null);

    await expect(
      service.update('missing', { isActive: false }),
    ).rejects.toBeInstanceOf(AppError);
    expect(repo.update).not.toHaveBeenCalled();
  });
});
