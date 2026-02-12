import { Test, TestingModule } from '@nestjs/testing';
import { WorkflowRunController } from './workflow-run.controller';
import { WorkflowRunService } from './workflow-run.service';
import { createWorkflowRunFixture } from 'test/workflow-run/workflow-run.fixtures';

describe('WorkflowRunController', () => {
  let controller: WorkflowRunController;
  let service: WorkflowRunService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [WorkflowRunController],
      providers: [
        {
          provide: WorkflowRunService,
          useValue: {
            list: jest.fn(),
            get: jest.fn(),
          },
        },
      ],
    }).compile();

    controller = module.get<WorkflowRunController>(WorkflowRunController);
    service = module.get<WorkflowRunService>(WorkflowRunService);
  });

  it('list() calls WorkflowRunService.list()', async () => {
    const list = [createWorkflowRunFixture({ id: 'wfr_1' })];
    const listSpy = jest.spyOn(service, 'list').mockResolvedValue({
      items: list,
      pagination: {
        page: 1,
        pageSize: 25,
        total: 1,
        totalPages: 1,
        hasNext: false,
        hasPrev: false,
      },
    });

    await expect(
      controller.list('wf_1', { page: 1, pageSize: 25 }),
    ).resolves.toEqual({
      items: list,
      pagination: {
        page: 1,
        pageSize: 25,
        total: 1,
        totalPages: 1,
        hasNext: false,
        hasPrev: false,
      },
    });
    expect(listSpy).toHaveBeenCalledWith({
      workflowId: 'wf_1',
      page: 1,
      pageSize: 25,
    });
  });

  it('get() calls WorkflowRunService.get()', async () => {
    const run = createWorkflowRunFixture({ id: 'wfr_1', workflowId: 'wf_1' });
    const getSpy = jest.spyOn(service, 'get').mockResolvedValue(run);

    await expect(controller.get('wf_1', 'wfr_1')).resolves.toBe(run);
    expect(getSpy).toHaveBeenCalledWith('wf_1', 'wfr_1');
  });
});
