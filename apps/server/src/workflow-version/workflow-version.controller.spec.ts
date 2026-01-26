import { Test, TestingModule } from '@nestjs/testing';
import { WorkflowVersionController } from './workflow-version.controller';
import { WorkflowVersionService } from './workflow-version.service';
import {
  createWorkflowVersionFixture,
  createWorkflowVersionListFixture,
} from 'test/workflow-version/workflow-version.fixtures';

describe('WorkflowVersionController', () => {
  let controller: WorkflowVersionController;
  let service: WorkflowVersionService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [WorkflowVersionController],
      providers: [
        {
          provide: WorkflowVersionService,
          useValue: {
            list: jest.fn(),
            get: jest.fn(),
          },
        },
      ],
    }).compile();

    controller = module.get<WorkflowVersionController>(
      WorkflowVersionController,
    );
    service = module.get<WorkflowVersionService>(WorkflowVersionService);
  });

  it('list() calls WorkflowVersionService.list()', async () => {
    const list = createWorkflowVersionListFixture(2);
    const listSpy = jest.spyOn(service, 'list').mockResolvedValue(list);

    await expect(controller.list('wf_1')).resolves.toBe(list);
    expect(listSpy).toHaveBeenCalledWith('wf_1');
  });

  it('get() calls WorkflowVersionService.get()', async () => {
    const version = createWorkflowVersionFixture({
      workflowId: 'wf_1',
      version: 2,
    });
    const getSpy = jest.spyOn(service, 'get').mockResolvedValue(version);

    await expect(controller.get('wf_1', 2)).resolves.toBe(version);
    expect(getSpy).toHaveBeenCalledWith('wf_1', 2);
  });
});
