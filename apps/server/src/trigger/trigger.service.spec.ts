import { Test } from '@nestjs/testing';
import { TriggerService } from './trigger.service';
import { TriggerRepository, WorkflowRepository } from '@taskforge/db-access';
import {
  createTriggerRepositoryMock,
  type TriggerRepositoryMock,
} from 'test/trigger/trigger.repository.mock';
import {
  createTriggerFixture,
  createTriggerListFixture,
} from 'test/trigger/trigger.fixtures';
import {
  createWorkflowRepositoryMock,
  type WorkflowRepositoryMock,
} from 'test/workflow/workflow.repository.mock';
import { createWorkflowFixture } from 'test/workflow/workflow.fixtures';
import { AppError } from 'src/common/http/errors/app-error';

describe('TriggerService', () => {
  let service: TriggerService;
  let repo: TriggerRepositoryMock;
  let workflowRepo: WorkflowRepositoryMock;

  beforeEach(async () => {
    repo = createTriggerRepositoryMock();
    workflowRepo = createWorkflowRepositoryMock();

    const moduleRef = await Test.createTestingModule({
      providers: [
        TriggerService,
        { provide: TriggerRepository, useValue: repo },
        { provide: WorkflowRepository, useValue: workflowRepo },
      ],
    }).compile();

    service = moduleRef.get(TriggerService);
  });

  it('create() creates a trigger for workflow', async () => {
    const wf = createWorkflowFixture({ id: 'wf_1' });
    const created = createTriggerFixture({
      workflowId: 'wf_1',
      type: 'WEBHOOK',
    });

    workflowRepo.findById.mockResolvedValue(wf);
    repo.create.mockResolvedValue(created);

    await expect(
      service.create({
        workflowId: 'wf_1',
        type: 'WEBHOOK',
        config: { url: 'https://example.com' },
      }),
    ).resolves.toBe(created);

    expect(repo.create).toHaveBeenCalledWith({
      workflow: { connect: { id: 'wf_1' } },
      type: 'WEBHOOK',
      name: undefined,
      isActive: true,
      config: { url: 'https://example.com' },
    });
  });

  it('create() throws notFound when workflow missing', async () => {
    workflowRepo.findById.mockResolvedValue(null);

    await expect(
      service.create({ workflowId: 'missing', type: 'MANUAL' }),
    ).rejects.toBeInstanceOf(AppError);
    expect(repo.create).not.toHaveBeenCalled();
  });

  it('list() returns triggers for workflow', async () => {
    const wf = createWorkflowFixture({ id: 'wf_1' });
    const list = createTriggerListFixture(2);

    workflowRepo.findById.mockResolvedValue(wf);
    repo.findPageByWorkflow.mockResolvedValue({ items: list, total: 2 });

    await expect(
      service.list({ workflowId: 'wf_1', page: 1, pageSize: 25 }),
    ).resolves.toEqual({
      items: list,
      pagination: {
        page: 1,
        pageSize: 25,
        total: 2,
        totalPages: 1,
        hasNext: false,
        hasPrev: false,
      },
    });
    expect(repo.findPageByWorkflow).toHaveBeenCalledWith({
      workflowId: 'wf_1',
      page: 1,
      pageSize: 25,
    });
  });

  it('get() returns trigger when found', async () => {
    const wf = createWorkflowFixture({ id: 'wf_1' });
    const trigger = createTriggerFixture({ id: 'tr_1', workflowId: 'wf_1' });

    workflowRepo.findById.mockResolvedValue(wf);
    repo.findById.mockResolvedValue(trigger);

    await expect(service.get('wf_1', 'tr_1')).resolves.toBe(trigger);
    expect(repo.findById).toHaveBeenCalledWith('tr_1');
  });

  it('get() throws notFound when trigger missing', async () => {
    const wf = createWorkflowFixture({ id: 'wf_1' });
    workflowRepo.findById.mockResolvedValue(wf);
    repo.findById.mockResolvedValue(null);

    await expect(service.get('wf_1', 'missing')).rejects.toBeInstanceOf(
      AppError,
    );
  });

  it('get() throws notFound when trigger belongs to another workflow', async () => {
    const wf = createWorkflowFixture({ id: 'wf_1' });
    const trigger = createTriggerFixture({ id: 'tr_1', workflowId: 'wf_2' });

    workflowRepo.findById.mockResolvedValue(wf);
    repo.findById.mockResolvedValue(trigger);

    await expect(service.get('wf_1', 'tr_1')).rejects.toBeInstanceOf(AppError);
  });

  it('delete() soft deletes trigger', async () => {
    const wf = createWorkflowFixture({ id: 'wf_1' });
    const trigger = createTriggerFixture({ id: 'tr_1', workflowId: 'wf_1' });
    const deleted = createTriggerFixture({
      id: 'tr_1',
      workflowId: 'wf_1',
      isActive: false,
    });

    workflowRepo.findById.mockResolvedValue(wf);
    repo.findById.mockResolvedValue(trigger);
    repo.softDelete.mockResolvedValue(deleted);

    await expect(service.delete('wf_1', 'tr_1')).resolves.toBe(deleted);
    expect(repo.softDelete).toHaveBeenCalledWith('tr_1');
  });

  it('update() updates trigger after existence check', async () => {
    const wf = createWorkflowFixture({ id: 'wf_1' });
    const trigger = createTriggerFixture({ id: 'tr_1', workflowId: 'wf_1' });
    const updated = createTriggerFixture({
      id: 'tr_1',
      workflowId: 'wf_1',
      name: 'Updated',
    });

    workflowRepo.findById.mockResolvedValue(wf);
    repo.findById.mockResolvedValue(trigger);
    repo.update.mockResolvedValue(updated);

    await expect(
      service.update('wf_1', 'tr_1', { name: 'Updated' }),
    ).resolves.toBe(updated);

    expect(repo.update).toHaveBeenCalledWith('tr_1', { name: 'Updated' });
  });
});
