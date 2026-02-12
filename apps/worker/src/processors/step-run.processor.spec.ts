import { StepRunProcessor } from './step-run.processor';
import {
  type PrismaService,
  SecretRepository,
  StepRunRepository,
  WorkflowRunRepository,
  WorkflowVersionRepository,
} from '@taskforge/db-access';
import type { Job } from 'bullmq';
import type { StepRunJobPayload } from '@taskforge/contracts';
import { ExecutorRegistry } from '../executors/executor-registry';
import { createPrismaServiceMock } from 'test/prisma.mocks';
import { CryptoService } from '../crypto/crypto.service';

describe('StepRunProcessor', () => {
  let executeMock: jest.Mock;
  let registryMock: jest.Mock;

  beforeEach(() => {
    jest.useFakeTimers();
    jest.clearAllMocks();

    executeMock = jest.fn().mockResolvedValue({
      statusCode: 200,
      headers: { 'content-type': 'application/json' },
      body: { ok: true },
    });

    registryMock = jest.fn().mockReturnValue({ stepType: 'http', execute: executeMock });
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it('marks step succeeded and workflow succeeded when all steps complete', async () => {
    const prisma = createPrismaServiceMock();

    const startedAt = new Date(Date.now() - 50);

    prisma.workflowRun.updateMany.mockResolvedValue({ count: 1 });
    prisma.workflowVersion.findUnique.mockResolvedValue({
      id: 'wfv_1',
      definition: {
        steps: [
          {
            key: 'step_1',
            type: 'http',
            request: {
              method: 'POST',
              url: 'https://x.test',
              query: { base: '1' },
              body: { content: 'base' },
            },
          },
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
    prisma.workflowRun.findUnique.mockResolvedValue({
      id: 'wfr_1',
      workflowId: 'wf_1',
    });

    const prismaService = prisma as unknown as PrismaService;
    const stepRunRepo = new StepRunRepository(prismaService);
    const workflowRunRepo = new WorkflowRunRepository(prismaService);
    const workflowVersionRepo = new WorkflowVersionRepository(prismaService);
    const secretRepo = {
      findManyByNames: jest.fn().mockResolvedValue([]),
    } as unknown as SecretRepository;

    process.env.TASKFORGE_SECRET_KEY = '0'.repeat(64);
    const crypto = new CryptoService();

    const redis = { eval: jest.fn().mockResolvedValue([1, 60]) };

    const processor = new StepRunProcessor(
      stepRunRepo,
      workflowRunRepo,
      workflowVersionRepo,
      secretRepo,
      crypto,
      redis as any,
      { get: registryMock } as unknown as ExecutorRegistry,
    );

    const job = {
      id: 'job_1',
      data: {
        workflowRunId: 'wfr_1',
        stepRunId: 'sr_1',
        stepKey: 'step_1',
        workflowVersionId: 'wfv_1',
        input: { hello: 'world' },
        requestOverride: {
          query: { override: '2' },
          body: { content: 'override' },
        },
      },
    } as unknown as Job<StepRunJobPayload>;

    const p = processor.process(job);
    await jest.runOnlyPendingTimersAsync();
    await p;

    expect(prisma.workflowRun.updateMany).toHaveBeenCalledWith(
      expect.objectContaining({ where: { id: 'wfr_1', status: 'QUEUED' } }),
    );

    expect(prisma.stepRun.update).toHaveBeenCalledWith(
      expect.objectContaining({ where: { id: 'sr_1' } }),
    );

    expect(executeMock).toHaveBeenCalledTimes(1);
    expect(executeMock).toHaveBeenCalledWith({
      request: {
        method: 'POST',
        url: 'https://x.test',
        query: { base: '1', override: '2' },
        body: { content: 'override' },
      },
      input: {
        input: { hello: 'world' },
        steps: {},
        secret: {},
      },
    });

    expect(prisma.workflowRun.update).toHaveBeenCalledWith(
      expect.objectContaining({
        where: { id: 'wfr_1' },
        data: expect.objectContaining({ status: 'SUCCEEDED' }) as unknown,
      }),
    );
  });

  it('continues when rate limit check fails', async () => {
    const prisma = createPrismaServiceMock();

    const startedAt = new Date(Date.now() - 50);

    prisma.workflowRun.updateMany.mockResolvedValue({ count: 1 });
    prisma.workflowVersion.findUnique.mockResolvedValue({
      id: 'wfv_1',
      definition: {
        steps: [
          {
            key: 'step_1',
            type: 'http',
            request: { method: 'GET', url: 'https://x.test' },
          },
        ],
      },
    });

    prisma.stepRun.update.mockResolvedValue({});
    prisma.stepRun.findUnique.mockResolvedValue({
      id: 'sr_1',
      workflowRunId: 'wfr_1',
      startedAt,
    });
    prisma.stepRun.findMany.mockResolvedValue([{ id: 'sr_1', status: 'SUCCEEDED' }]);

    prisma.workflowRun.update.mockResolvedValue({});
    prisma.workflowRun.findUnique.mockResolvedValue({
      id: 'wfr_1',
      workflowId: 'wf_1',
    });

    const prismaService = prisma as unknown as PrismaService;
    const stepRunRepo = new StepRunRepository(prismaService);
    const workflowRunRepo = new WorkflowRunRepository(prismaService);
    const workflowVersionRepo = new WorkflowVersionRepository(prismaService);
    const secretRepo = {
      findManyByNames: jest.fn().mockResolvedValue([]),
    } as unknown as SecretRepository;

    process.env.TASKFORGE_SECRET_KEY = '0'.repeat(64);
    const crypto = new CryptoService();

    const redis = { eval: jest.fn().mockRejectedValue(new Error('redis down')) };

    const processor = new StepRunProcessor(
      stepRunRepo,
      workflowRunRepo,
      workflowVersionRepo,
      secretRepo,
      crypto,
      redis as any,
      { get: registryMock } as unknown as ExecutorRegistry,
    );

    const job = {
      id: 'job_1',
      data: {
        workflowRunId: 'wfr_1',
        stepRunId: 'sr_1',
        stepKey: 'step_1',
        workflowVersionId: 'wfv_1',
        input: { hello: 'world' },
      },
    } as unknown as Job<StepRunJobPayload>;

    const p = processor.process(job);
    await jest.runOnlyPendingTimersAsync();
    await p;

    expect(executeMock).toHaveBeenCalledTimes(1);
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
    prisma.workflowRun.findUnique.mockResolvedValue({
      id: 'wfr_1',
      workflowId: 'wf_1',
    });

    const prismaService = prisma as unknown as PrismaService;
    const stepRunRepo = new StepRunRepository(prismaService);
    const workflowRunRepo = new WorkflowRunRepository(prismaService);
    const workflowVersionRepo = new WorkflowVersionRepository(prismaService);
    const secretRepo = {
      findManyByNames: jest.fn().mockResolvedValue([]),
    } as unknown as SecretRepository;

    process.env.TASKFORGE_SECRET_KEY = '0'.repeat(64);
    const crypto = new CryptoService();

    const executeFailed = jest.fn().mockRejectedValue(new Error('boom'));

    const redis = { eval: jest.fn().mockResolvedValue([1, 60]) };

    const processor = new StepRunProcessor(
      stepRunRepo,
      workflowRunRepo,
      workflowVersionRepo,
      secretRepo,
      crypto,
      redis as any,
      { get: () => ({ stepType: 'http', execute: executeFailed }) } as unknown as ExecutorRegistry,
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
