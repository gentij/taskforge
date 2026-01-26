import { Test, TestingModule } from '@nestjs/testing';
import { WorkflowController } from './workflow.controller';
import { WorkflowService } from './workflow.service';
import { createWorkflowFixture } from 'test/workflow/workflow.fixtures';
import { createWorkflowVersionFixture } from 'test/workflow-version/workflow-version.fixtures';

describe('WorkflowController', () => {
  let controller: WorkflowController;
  let service: WorkflowService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [WorkflowController],
      providers: [
        {
          provide: WorkflowService,
          useValue: {
            create: jest.fn(),
            list: jest.fn(),
            get: jest.fn(),
            update: jest.fn(),
            createVersion: jest.fn(),
          },
        },
      ],
    }).compile();

    controller = module.get<WorkflowController>(WorkflowController);
    service = module.get<WorkflowService>(WorkflowService);
  });

  it('create() calls WorkflowService.create()', async () => {
    const wf = createWorkflowFixture({ name: 'My WF' });
    const createSpy = jest.spyOn(service, 'create').mockResolvedValue(wf);

    await expect(controller.create({ name: 'My WF' })).resolves.toBe(wf);
    expect(createSpy).toHaveBeenCalledWith('My WF');
  });

  it('list() calls WorkflowService.list()', async () => {
    const list = [createWorkflowFixture({ id: 'wf_1' })];
    const listSpy = jest.spyOn(service, 'list').mockResolvedValue(list);

    await expect(controller.list()).resolves.toBe(list);
    expect(listSpy).toHaveBeenCalledTimes(1);
  });

  it('get() calls WorkflowService.get()', async () => {
    const wf = createWorkflowFixture({ id: 'wf_1' });
    const getSpy = jest.spyOn(service, 'get').mockResolvedValue(wf);

    await expect(controller.get('wf_1')).resolves.toBe(wf);
    expect(getSpy).toHaveBeenCalledWith('wf_1');
  });

  it('update() calls WorkflowService.update()', async () => {
    const wf = createWorkflowFixture({ id: 'wf_1', name: 'New' });
    const updateSpy = jest.spyOn(service, 'update').mockResolvedValue(wf);

    await expect(controller.update('wf_1', { name: 'New' })).resolves.toBe(wf);
    expect(updateSpy).toHaveBeenCalledWith('wf_1', { name: 'New' });
  });

  it('createVersion() calls WorkflowService.createVersion()', async () => {
    const version = createWorkflowVersionFixture({
      workflowId: 'wf_1',
      version: 2,
    });
    const createVersionSpy = jest
      .spyOn(service, 'createVersion')
      .mockResolvedValue(version);

    await expect(
      controller.createVersion('wf_1', { definition: { steps: [] } }),
    ).resolves.toBe(version);
    expect(createVersionSpy).toHaveBeenCalledWith('wf_1', {
      steps: [],
    });
  });
});
