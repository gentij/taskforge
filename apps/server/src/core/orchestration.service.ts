import { Injectable, Logger } from '@nestjs/common';
import { Prisma } from '@prisma/client';
import { PrismaService } from '@taskforge/db-access';
import { StepRunJobPayload, WorkflowDefinition } from '@taskforge/contracts';

import { StepRunQueueService } from '../queue/step-run-queue.service';

type EnqueueItem = {
  stepType: string;
  payload: StepRunJobPayload;
};

@Injectable()
export class OrchestrationService {
  private readonly logger = new Logger(OrchestrationService.name);

  constructor(
    private readonly prisma: PrismaService,
    private readonly stepRunQueueService: StepRunQueueService,
  ) {}

  async startWorkflow(params: {
    workflowId: string;
    workflowVersionId: string;
    triggerId?: string;
    eventType?: string;
    eventExternalId?: string;
    eventPayload?: Record<string, unknown>;
    input?: Record<string, unknown>;
    overrides?: Record<
      string,
      {
        query?: Record<string, string | number | boolean>;
        body?: unknown;
      }
    >;
  }): Promise<{ workflowRunId: string; stepRunIds: string[] }> {
    const {
      workflowId,
      workflowVersionId,
      triggerId,
      eventType,
      eventExternalId,
      eventPayload,
      input,
      overrides,
    } = params;

    this.logger.log(
      `Starting workflow workflowId=${workflowId} workflowVersionId=${workflowVersionId} triggerId=${triggerId ?? '-'} eventType=${eventType ?? '-'}`,
    );

    const { workflowRunId, stepRunIds, enqueueItems } =
      await this.prisma.$transaction(async (tx) => {
        const triggerIdToUse = await this.resolveTriggerIdForEvent(tx, {
          workflowId,
          triggerId,
        });

        const event = await tx.event.create({
          data: {
            triggerId: triggerIdToUse,
            type: eventType ?? 'MANUAL',
            externalId: eventExternalId,
            payload: (eventPayload ??
              input ??
              {}) as unknown as Prisma.InputJsonValue,
            receivedAt: new Date(),
          },
        });

        const workflowVersion = await tx.workflowVersion.findUniqueOrThrow({
          where: { id: workflowVersionId },
        });

        const definition =
          workflowVersion.definition as unknown as WorkflowDefinition;
        const steps = definition.steps ?? [];

        const overridesToPersist = this.pickAndNormalizeOverridesForSteps(
          steps,
          overrides,
        );

        const workflowRunData: Prisma.WorkflowRunUncheckedCreateInput = {
          workflowId,
          workflowVersionId,
          triggerId: triggerIdToUse,
          eventId: event.id,
          status: 'QUEUED',
          input: (input ?? {}) as unknown as Prisma.InputJsonValue,
          overrides: overridesToPersist
            ? (overridesToPersist as unknown as Prisma.InputJsonValue)
            : Prisma.DbNull,
        };

        const workflowRun = await tx.workflowRun.create({
          data: workflowRunData,
        });

        if (steps.length === 0) {
          await tx.workflowRun.update({
            where: { id: workflowRun.id },
            data: { status: 'SUCCEEDED', finishedAt: new Date() },
          });

          return {
            workflowRunId: workflowRun.id,
            stepRunIds: [] as string[],
            enqueueItems: [] as EnqueueItem[],
          };
        }

        const createdStepRuns = await Promise.all(
          steps.map((step) => {
            const stepRunData: Prisma.StepRunUncheckedCreateInput = {
              workflowRunId: workflowRun.id,
              stepKey: step.key,
              status: 'QUEUED',
              attempt: 0,
              input: (input ?? {}) as unknown as Prisma.InputJsonValue,
              requestOverride: overridesToPersist?.[step.key]
                ? (overridesToPersist[
                    step.key
                  ] as unknown as Prisma.InputJsonValue)
                : Prisma.DbNull,
            };

            return tx.stepRun.create({
              data: stepRunData,
              select: { id: true, stepKey: true },
            });
          }),
        );

        const stepRunIds = createdStepRuns.map((sr) => sr.id);

        const enqueueItems: EnqueueItem[] = createdStepRuns.map((sr, idx) => {
          const step = steps[idx];
          return {
            stepType: step.type,
            payload: {
              workflowRunId: workflowRun.id,
              stepRunId: sr.id,
              stepKey: sr.stepKey,
              workflowVersionId,
              input: input ?? {},
              requestOverride: overridesToPersist?.[step.key],
            },
          };
        });

        return { workflowRunId: workflowRun.id, stepRunIds, enqueueItems };
      });

    // Enqueue after DB transaction commits
    const enqueuedStepRunIds: string[] = [];
    try {
      for (const item of enqueueItems) {
        await this.stepRunQueueService.enqueueStepRun(
          item.stepType,
          item.payload,
        );
        enqueuedStepRunIds.push(item.payload.stepRunId);
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);

      this.logger.error(
        `Failed to enqueue step runs for workflowRunId=${workflowRunId}: ${message} (enqueued ${enqueuedStepRunIds.length}/${enqueueItems.length})`,
      );

      const notEnqueued = stepRunIds.filter(
        (id) => !enqueuedStepRunIds.includes(id),
      );

      await this.prisma.$transaction(async (tx) => {
        if (notEnqueued.length > 0) {
          await tx.stepRun.updateMany({
            where: { id: { in: notEnqueued }, status: 'QUEUED' },
            data: {
              status: 'FAILED',
              finishedAt: new Date(),
              error: {
                message: 'Failed to enqueue step run',
                details: message,
              } as Prisma.InputJsonValue,
            },
          });
        }

        await tx.workflowRun.update({
          where: { id: workflowRunId },
          data: {
            status: 'FAILED',
            finishedAt: new Date(),
          },
        });
      });

      throw err;
    }

    this.logger.log(
      `Workflow started workflowRunId=${workflowRunId} steps=${stepRunIds.length}`,
    );

    return { workflowRunId, stepRunIds };
  }

  private async resolveTriggerIdForEvent(
    tx: Prisma.TransactionClient,
    params: {
      workflowId: string;
      triggerId?: string;
    },
  ): Promise<string> {
    if (params.triggerId) return params.triggerId;

    const manual = await tx.trigger.findFirst({
      where: { workflowId: params.workflowId, type: 'MANUAL' },
      select: { id: true },
    });

    if (manual) return manual.id;

    const created = await tx.trigger.create({
      data: {
        workflowId: params.workflowId,
        type: 'MANUAL',
        name: 'Manual',
        isActive: true,
        config: {},
      },
      select: { id: true },
    });

    return created.id;
  }

  private pickAndNormalizeOverridesForSteps(
    steps: Array<{ key: string }>,
    overrides:
      | Record<
          string,
          { query?: Record<string, string | number | boolean>; body?: unknown }
        >
      | undefined,
  ): Record<
    string,
    { query?: Record<string, string | number | boolean>; body?: unknown }
  > | null {
    if (!overrides) return null;

    const picked: Record<
      string,
      { query?: Record<string, string | number | boolean>; body?: unknown }
    > = {};

    for (const step of steps) {
      const raw = overrides[step.key];
      if (!raw) continue;

      const normalized: {
        query?: Record<string, string | number | boolean>;
        body?: unknown;
      } = {};

      if (raw.query && Object.keys(raw.query).length > 0) {
        normalized.query = raw.query;
      }

      if (raw.body !== undefined) {
        normalized.body = raw.body;
      }

      if (normalized.query !== undefined || normalized.body !== undefined) {
        picked[step.key] = normalized;
      }
    }

    return Object.keys(picked).length > 0 ? picked : null;
  }
}
