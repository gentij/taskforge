/* eslint-disable
  @typescript-eslint/no-unsafe-assignment,
  @typescript-eslint/no-unsafe-member-access
*/

import { Test } from '@nestjs/testing';
import { APP_FILTER, APP_GUARD, APP_INTERCEPTOR, APP_PIPE } from '@nestjs/core';
import {
  FastifyAdapter,
  type NestFastifyApplication,
} from '@nestjs/platform-fastify';
import { ZodSerializerInterceptor, ZodValidationPipe } from 'nestjs-zod';

import { WorkflowController } from 'src/workflow/workflow.controller';
import { WorkflowService } from 'src/workflow/workflow.service';
import { WorkflowRepository } from 'src/workflow/workflow.repository';
import { PrismaService } from 'src/prisma/prisma.service';

import { AllExceptionsFilter } from 'src/common/http/filters/all-exceptions.filter'; // <-- adjust path to your file

import {
  createWorkflowRepositoryMock,
  type WorkflowRepositoryMock,
} from 'test/workflow/workflow.repository.mock';
import {
  createWorkflowFixture,
  createWorkflowListFixture,
} from 'test/workflow/workflow.fixtures';
import { createWorkflowVersionFixture } from 'test/workflow-version/workflow-version.fixtures';
import { AllowAuthGuard } from 'test/utils/allow-auth.guard';
import { ResponseInterceptor } from 'src/common/http/interceptors/response.interceptor';

describe('Workflow (e2e)', () => {
  let app: NestFastifyApplication;
  let repo: WorkflowRepositoryMock;
  let prisma: {
    $transaction: jest.Mock<Promise<unknown>, [(tx: any) => unknown]>;
  };

  beforeEach(async () => {
    repo = createWorkflowRepositoryMock();
    prisma = {
      $transaction: jest.fn<Promise<unknown>, [(tx: any) => unknown]>(),
    };

    const moduleRef = await Test.createTestingModule({
      controllers: [WorkflowController],
      providers: [
        WorkflowService,
        { provide: WorkflowRepository, useValue: repo },
        { provide: PrismaService, useValue: prisma },

        { provide: APP_PIPE, useClass: ZodValidationPipe },
        { provide: APP_INTERCEPTOR, useClass: ZodSerializerInterceptor },
        { provide: APP_FILTER, useClass: AllExceptionsFilter },

        { provide: APP_GUARD, useClass: AllowAuthGuard },
        { provide: APP_INTERCEPTOR, useClass: ResponseInterceptor },
      ],
    }).compile();

    app = moduleRef.createNestApplication<NestFastifyApplication>(
      new FastifyAdapter(),
    );

    await app.init();
    await app.getHttpAdapter().getInstance().ready();
  });

  afterEach(async () => {
    await app.close();
  });

  it('POST /workflows -> 201 + ok:true + data', async () => {
    const created = createWorkflowFixture({ id: 'wf_new', name: 'My WF' });
    const version = createWorkflowVersionFixture({
      workflowId: created.id,
      version: 1,
      definition: { steps: [] },
    });
    const updated = createWorkflowFixture({
      id: created.id,
      name: 'My WF',
      latestVersionId: version.id,
    });

    const tx = {
      workflow: {
        create: jest.fn().mockResolvedValue(created),
        update: jest.fn().mockResolvedValue(updated),
      },
      workflowVersion: {
        create: jest.fn().mockResolvedValue(version),
        findFirst: jest.fn(),
      },
    };

    prisma.$transaction.mockImplementation((cb) => Promise.resolve(cb(tx)));

    const res = await app.inject({
      method: 'POST',
      url: '/workflows',
      payload: { name: 'My WF' },
    });

    expect(res.statusCode).toBe(201);

    const body = res.json();

    expect(body.ok).toBe(true);
    expect(body.data.name).toBe('My WF');
    expect(tx.workflow.create).toHaveBeenCalledWith({
      data: { name: 'My WF' },
    });
  });

  it('GET /workflows -> 200 + data array', async () => {
    const list = createWorkflowListFixture(2);
    repo.findMany.mockResolvedValue(list);

    const res = await app.inject({
      method: 'GET',
      url: '/workflows',
    });

    expect(res.statusCode).toBe(200);

    const body = res.json();
    expect(body.ok).toBe(true);
    expect(Array.isArray(body.data)).toBe(true);
    expect(body.data).toHaveLength(2);

    expect(repo.findMany).toHaveBeenCalledTimes(1);
  });

  it('GET /workflows/:id -> 200 when found', async () => {
    const wf = createWorkflowFixture({ id: 'wf_1' });
    repo.findById.mockResolvedValue(wf);

    const res = await app.inject({
      method: 'GET',
      url: '/workflows/wf_1',
    });

    expect(res.statusCode).toBe(200);

    const body = res.json();
    expect(body.ok).toBe(true);
    expect(body.data.id).toBe('wf_1');
  });

  it('GET /workflows/:id -> 404 with standardized error payload', async () => {
    repo.findById.mockResolvedValue(null);

    const res = await app.inject({
      method: 'GET',
      url: '/workflows/missing',
    });

    expect(res.statusCode).toBe(404);

    const body = res.json();
    expect(body.ok).toBe(false);
    expect(body.error).toBeDefined();
    expect(typeof body.error.code).toBe('string');
  });

  it('PATCH /workflows/:id -> 200 updates workflow', async () => {
    const existing = createWorkflowFixture({ id: 'wf_1', name: 'Old' });
    const updated = createWorkflowFixture({ id: 'wf_1', name: 'New' });

    repo.findById.mockResolvedValue(existing);
    repo.update.mockResolvedValue(updated);

    const res = await app.inject({
      method: 'PATCH',
      url: '/workflows/wf_1',
      payload: { name: 'New' },
    });

    expect(res.statusCode).toBe(200);

    const body = res.json();
    expect(body.ok).toBe(true);
    expect(body.data.name).toBe('New');

    expect(repo.update).toHaveBeenCalledWith('wf_1', { name: 'New' });
  });

  it('POST /workflows/:id/versions -> 201 creates a new version', async () => {
    const wf = createWorkflowFixture({ id: 'wf_1' });
    const version = createWorkflowVersionFixture({
      id: 'wfv_2',
      workflowId: 'wf_1',
      version: 2,
      definition: { steps: [{ id: 's1' }] },
    });

    repo.findById.mockResolvedValue(wf);

    const tx = {
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

    const res = await app.inject({
      method: 'POST',
      url: '/workflows/wf_1/versions',
      payload: { definition: version.definition },
    });

    expect(res.statusCode).toBe(201);

    const body = res.json();
    expect(body.ok).toBe(true);
    expect(body.data.version).toBe(2);
    expect(body.data.workflowId).toBe('wf_1');
  });

  it('POST /workflows -> 400 validation error when name missing', async () => {
    const res = await app.inject({
      method: 'POST',
      url: '/workflows',
      payload: {}, // invalid
    });

    expect(res.statusCode).toBe(400);

    const body = res.json();
    expect(body.ok).toBe(false);
    expect(body.error.code).toBe('VALIDATION_ERROR');
    expect(Array.isArray(body.error.details)).toBe(true);
  });
});
