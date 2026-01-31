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

import { TriggerController } from 'src/trigger/trigger.controller';
import { TriggerService } from 'src/trigger/trigger.service';
import { OrchestrationService } from 'src/core/orchestration.service';
import { WorkflowService } from 'src/workflow/workflow.service';
import { AllExceptionsFilter } from 'src/common/http/filters/all-exceptions.filter';
import { ResponseInterceptor } from 'src/common/http/interceptors/response.interceptor';
import { AllowAuthGuard } from 'test/utils/allow-auth.guard';

import {
  createTriggerRepositoryMock,
  type TriggerRepositoryMock,
} from 'test/trigger/trigger.repository.mock';
import {
  createWorkflowRepositoryMock,
  type WorkflowRepositoryMock,
} from 'test/workflow/workflow.repository.mock';
import { TriggerRepository, WorkflowRepository } from '@taskforge/db-access';
import { createTriggerFixture } from 'test/trigger/trigger.fixtures';
import { createWorkflowFixture } from 'test/workflow/workflow.fixtures';

describe('Trigger Webhook (integration e2e)', () => {
  let app: NestFastifyApplication;
  let triggerRepo: TriggerRepositoryMock;
  let workflowRepo: WorkflowRepositoryMock;
  let orchestration: { startWorkflow: jest.Mock };
  let workflowService: { get: jest.Mock };

  beforeEach(async () => {
    triggerRepo = createTriggerRepositoryMock();
    workflowRepo = createWorkflowRepositoryMock();
    orchestration = { startWorkflow: jest.fn() };
    workflowService = { get: jest.fn() };

    const moduleRef = await Test.createTestingModule({
      controllers: [TriggerController],
      providers: [
        TriggerService,
        { provide: TriggerRepository, useValue: triggerRepo },
        { provide: WorkflowRepository, useValue: workflowRepo },
        { provide: OrchestrationService, useValue: orchestration },
        { provide: WorkflowService, useValue: workflowService },

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
    await app?.close();
  });

  it('POST /workflows/:workflowId/triggers/:id/webhook enqueues a workflow run', async () => {
    workflowRepo.findById.mockResolvedValue(
      createWorkflowFixture({ id: 'wf_1', latestVersionId: 'wfv_1' }),
    );
    triggerRepo.findById.mockResolvedValue(
      createTriggerFixture({ id: 'tr_1', workflowId: 'wf_1', isActive: true }),
    );
    workflowService.get.mockResolvedValue(
      createWorkflowFixture({ id: 'wf_1', latestVersionId: 'wfv_1' }),
    );
    orchestration.startWorkflow.mockResolvedValue({
      workflowRunId: 'wfr_1',
      stepRunIds: [],
    });

    const res = await app.inject({
      method: 'POST',
      url: '/workflows/wf_1/triggers/tr_1/webhook',
      payload: { hello: 'world' },
    });

    expect(res.statusCode).toBe(201);
    const body = res.json();
    expect(body.ok).toBe(true);
    expect(body.data.status).toBe('accepted');

    expect(orchestration.startWorkflow).toHaveBeenCalledWith(
      expect.objectContaining({
        workflowId: 'wf_1',
        workflowVersionId: 'wfv_1',
        triggerId: 'tr_1',
        eventType: 'WEBHOOK',
      }),
    );
  });
});
