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
import { createWorkflowVersionFixture } from 'test/workflow-version/workflow-version.fixtures';
import { PrismaService } from 'src/prisma/prisma.service';
import { AppError } from 'src/common/http/errors/app-error';
import {
  createPrismaServiceMock,
  PrismaServiceMock,
  PrismaTxMock,
} from 'test/prisma/prisma.mocks';

describe('WorkflowService', () => {
  let service: WorkflowService;
  let repo: WorkflowRepositoryMock;
  let prisma: PrismaServiceMock;

  beforeEach(async () => {
    repo = createWorkflowRepositoryMock();
    prisma = createPrismaServiceMock();

    const moduleRef = await Test.createTestingModule({
      providers: [
        WorkflowService,
        { provide: WorkflowRepository, useValue: repo },
        { provide: PrismaService, useValue: prisma },
      ],
    }).compile();

    service = moduleRef.get(WorkflowService);
  });

  it('create() creates a workflow', async () => {
    const wf = createWorkflowFixture({ name: 'My WF' });
    const version = createWorkflowVersionFixture({
      workflowId: wf.id,
      version: 1,
      definition: { steps: [] },
    });
    const updated = createWorkflowFixture({
      id: wf.id,
      name: 'My WF',
      latestVersionId: version.id,
    });

    const tx: PrismaTxMock = {
      workflow: {
        create: jest.fn().mockResolvedValue(wf),
        update: jest.fn().mockResolvedValue(updated),
      },
      workflowVersion: {
        create: jest.fn().mockResolvedValue(version),
        findFirst: jest.fn(),
      },
    };

    prisma.$transaction.mockImplementation((cb) => Promise.resolve(cb(tx)));

    await expect(service.create('My WF')).resolves.toBe(updated);

    expect(prisma.$transaction).toHaveBeenCalledTimes(1);
    expect(tx.workflow.create).toHaveBeenCalledWith({
      data: { name: 'My WF' },
    });
    expect(tx.workflowVersion.create).toHaveBeenCalledWith({
      data: {
        workflowId: wf.id,
        version: 1,
        definition: { steps: [] },
      },
    });
    expect(tx.workflow.update).toHaveBeenCalledWith({
      where: { id: wf.id },
      data: { latestVersionId: version.id },
    });
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

  it('createVersion() creates next version and updates latestVersionId', async () => {
    const wf = createWorkflowFixture({ id: 'wf_1' });
    repo.findById.mockResolvedValue(wf);

    const version = createWorkflowVersionFixture({
      id: 'wfv_2',
      workflowId: 'wf_1',
      version: 2,
      definition: { steps: [{ id: 'step-1' }] },
    });

    const tx: PrismaTxMock = {
      workflow: {
        create: jest.fn(),
        update: jest.fn().mockResolvedValue(wf),
      },
      workflowVersion: {
        findFirst: jest.fn().mockResolvedValue({ version: 1 }),
        create: jest.fn().mockResolvedValue(version),
      },
    };

    prisma.$transaction.mockImplementation((cb) => Promise.resolve(cb(tx)));

    await expect(
      service.createVersion('wf_1', version.definition),
    ).resolves.toBe(version);

    expect(tx.workflowVersion.create).toHaveBeenCalledWith({
      data: {
        workflowId: 'wf_1',
        version: 2,
        definition: version.definition,
      },
    });
    expect(tx.workflow.update).toHaveBeenCalledWith({
      where: { id: 'wf_1' },
      data: { latestVersionId: version.id },
    });
  });

  it('createVersion() throws notFound when workflow missing', async () => {
    repo.findById.mockResolvedValue(null);

    await expect(
      service.createVersion('missing', { steps: [] }),
    ).rejects.toBeInstanceOf(AppError);
    expect(prisma.$transaction).not.toHaveBeenCalled();
  });
});
