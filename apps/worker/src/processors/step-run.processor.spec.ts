import { StepRunProcessor } from './step-run.processor';
import {
  type PrismaService,
  StepRunRepository,
  WorkflowRunRepository,
  WorkflowVersionRepository,
} from '@taskforge/db-access';
import type { Job } from 'bullmq';
import type { StepRunJobPayload } from '@taskforge/contracts';
import { ExecutorRegistry } from '../executors/executor-registry';
import { createPrismaServiceMock } from 'test/prisma.mocks';

describe('StepRunProcessor', () => {
  it('marks step succeeded and workflow succeeded when all steps complete', async () => {
    const prisma = createPrismaServiceMock();

    const startedAt = new Date(Date.now() - 50);

    prisma.workflowRun.updateMany.mockResolvedValue({ count: 1 });
    prisma.workflowVersion.findUnique.mockResolvedValue({
      id: 'wfv_1',
      definition: {
        steps: [
          { key: 'step_1', type: 'http', request: { method: 'GET', url: 'https://x.test' } },
          { key: 'step_2', type: 'http', request: { method: 'GET', url: 'https://x.test' } },
        ],
      },
    });

    prisma.stepRun.update.mockResolvedValue({});
    prisma.stepRun.findUnique.mockResolvedValue({
      id: 'sr_1',
      workflowRunId: 'wfr_1',
      startedAt,
    });
    prisma.stepRun.findMany.mockResolvedValue([
      { id: 'sr_1', status: 'SUCCEEDED' },
      { id: 'sr_2', status: 'SUCCEEDED' },
    ]);

    prisma.workflowRun.update.mockResolvedValue({});

    const prismaService = prisma as unknown as PrismaService;
    const stepRunRepo = new StepRunRepository(prismaService);
    const workflowRunRepo = new WorkflowRunRepository(prismaService);
    const workflowVersionRepo = new WorkflowVersionRepository(prismaService);

    const execute = jest.fn().mockResolvedValue({
      statusCode: 200,
      headers: { 'content-type': 'application/json' },
      body: { ok: true },
    });
    const registry = {
      get: jest.fn().mockReturnValue({ stepType: 'http', execute }),
    } as unknown as ExecutorRegistry;

    const processor = new StepRunProcessor(
      stepRunRepo,
      workflowRunRepo,
      workflowVersionRepo,
      registry,
    );

    const job = {
      id: 'job_1',
      data: {
        workflowRunId: 'wfr_1',
        stepRunId: 'sr_1',
        stepKey: 'step_1',
        workflowVersionId: 'wfv_1',
        input: {},
      },
    } as unknown as Job<StepRunJobPayload>;

    await processor.process(job);

    expect(prisma.workflowRun.updateMany).toHaveBeenCalledWith(
      expect.objectContaining({ where: { id: 'wfr_1', status: 'QUEUED' } }),
    );

    expect(prisma.stepRun.update).toHaveBeenCalledWith(
      expect.objectContaining({ where: { id: 'sr_1' } }),
    );

    expect(execute).toHaveBeenCalledTimes(1);

    expect(prisma.workflowRun.update).toHaveBeenCalledWith(
      expect.objectContaining({
        where: { id: 'wfr_1' },
        data: expect.objectContaining({ status: 'SUCCEEDED' }) as unknown,
      }),
    );
  });

  it('marks workflow failed when a sibling fails and all steps complete', async () => {
    const prisma = createPrismaServiceMock();

    prisma.workflowRun.updateMany.mockResolvedValue({ count: 1 });
    prisma.workflowVersion.findUnique.mockResolvedValue({
      id: 'wfv_1',
      definition: {
        steps: [{ key: 'step_1', type: 'http', request: { method: 'GET', url: 'https://x.test' } }],
      },
    });

    prisma.stepRun.update.mockResolvedValue({});
    prisma.stepRun.findUnique.mockResolvedValue({
      id: 'sr_1',
      workflowRunId: 'wfr_1',
      startedAt: new Date(),
    });

    prisma.stepRun.findMany.mockResolvedValue([
      { id: 'sr_1', status: 'FAILED' },
      { id: 'sr_2', status: 'SUCCEEDED' },
    ]);

    prisma.workflowRun.update.mockResolvedValue({});

    const prismaService = prisma as unknown as PrismaService;
    const stepRunRepo = new StepRunRepository(prismaService);
    const workflowRunRepo = new WorkflowRunRepository(prismaService);
    const workflowVersionRepo = new WorkflowVersionRepository(prismaService);

    const execute = jest.fn().mockRejectedValue(new Error('boom'));
    const registry = {
      get: jest.fn().mockReturnValue({ stepType: 'http', execute }),
    } as unknown as ExecutorRegistry;

    const processor = new StepRunProcessor(
      stepRunRepo,
      workflowRunRepo,
      workflowVersionRepo,
      registry,
    );

    const job = {
      id: 'job_1',
      data: {
        workflowRunId: 'wfr_1',
        stepRunId: 'sr_1',
        stepKey: 'step_1',
        workflowVersionId: 'wfv_1',
        input: {},
      },
    } as unknown as Job<StepRunJobPayload>;

    await expect(processor.process(job)).rejects.toThrow('boom');

    expect(prisma.workflowRun.update).toHaveBeenCalledWith(
      expect.objectContaining({
        where: { id: 'wfr_1' },
        data: expect.objectContaining({ status: 'FAILED' }) as unknown,
      }),
    );
  });
});
