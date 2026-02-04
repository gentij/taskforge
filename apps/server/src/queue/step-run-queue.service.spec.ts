import { Test } from '@nestjs/testing';
import { getQueueToken } from '@nestjs/bullmq';
import { StepRunQueueService } from './step-run-queue.service';
import { STEP_RUN_QUEUE_NAME } from './queue.constants';

describe('StepRunQueueService', () => {
  let service: StepRunQueueService;
  const queue = { add: jest.fn() };

  beforeEach(async () => {
    const moduleRef = await Test.createTestingModule({
      providers: [
        StepRunQueueService,
        { provide: getQueueToken(STEP_RUN_QUEUE_NAME), useValue: queue },
      ],
    }).compile();

    service = moduleRef.get(StepRunQueueService);
    queue.add.mockReset();
  });

  it('enqueueHttpStepRun() enqueues http job with defaults', async () => {
    const payload = {
      workflowRunId: 'wfr_1',
      stepRunId: 'sr_1',
      stepKey: 'step-1',
      workflowVersionId: 'wfv_1',
      input: { foo: 'bar' },
      dependsOn: [],
    };

    queue.add.mockResolvedValue({ id: 'sr_1' });

    await service.enqueueHttpStepRun(payload);

    expect(queue.add).toHaveBeenCalledWith(
      'http',
      payload,
      expect.objectContaining({
        jobId: 'sr_1',
        attempts: 3,
        backoff: { type: 'exponential', delay: 5000 },
        removeOnComplete: 1000,
        removeOnFail: 1000,
      }),
    );
  });

  it('enqueueHttpStepRun() merges job options', async () => {
    const payload = {
      workflowRunId: 'wfr_2',
      stepRunId: 'sr_2',
      stepKey: 'step-2',
      workflowVersionId: 'wfv_2',
      input: {},
      dependsOn: [],
    };

    queue.add.mockResolvedValue({ id: 'sr_2' });

    await service.enqueueHttpStepRun(payload, { attempts: 1 });

    expect(queue.add).toHaveBeenCalledWith(
      'http',
      payload,
      expect.objectContaining({
        attempts: 1,
        jobId: 'sr_2',
      }),
    );
  });
});
