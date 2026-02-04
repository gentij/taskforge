import { Processor, WorkerHost } from '@nestjs/bullmq';
import { Logger } from '@nestjs/common';
import { Job } from 'bullmq';
import {
  StepRunRepository,
  WorkflowRunRepository,
  WorkflowVersionRepository,
} from '@taskforge/db-access';
import { StepRunJobPayload } from '@taskforge/contracts';
import { ExecutorRegistry } from '../executors/executor-registry';
import { TemplateResolver } from '../utils/template-resolver';

interface StepDefinition {
  key: string;
  type: string;
  dependsOn?: string[];
  input?: Record<string, unknown>;
  request?: Record<string, unknown>;
}

@Processor('step-runs')
export class StepRunProcessor extends WorkerHost {
  private readonly logger = new Logger(StepRunProcessor.name);
  private readonly templateResolver = new TemplateResolver();

  constructor(
    private readonly stepRunRepository: StepRunRepository,
    private readonly workflowRunRepository: WorkflowRunRepository,
    private readonly workflowVersionRepository: WorkflowVersionRepository,
    private readonly executorRegistry: ExecutorRegistry,
  ) {
    super();
  }

  async process(job: Job<StepRunJobPayload>): Promise<void> {
    const { workflowRunId, stepRunId, stepKey, workflowVersionId, dependsOn, input: triggerInput } = job.data;

    try {
      await this.workflowRunRepository.markRunningIfQueued(workflowRunId);

      await this.stepRunRepository.update(stepRunId, {
        status: 'RUNNING',
        startedAt: new Date(),
      });

      const workflowVersion = await this.workflowVersionRepository.findById(workflowVersionId);

      if (!workflowVersion) {
        throw new Error(`WorkflowVersion not found: ${workflowVersionId}`);
      }

      const definition = workflowVersion.definition as unknown as {
        input?: Record<string, unknown>;
        steps: StepDefinition[];
      };
      const stepDef = definition.steps.find((s) => s.key === stepKey);

      if (!stepDef) {
        throw new Error(`Step definition not found for key: ${stepKey}`);
      }

      const mergedInput = (triggerInput as Record<string, unknown>) ?? {};
      const stepLevelInput = (stepDef.input ?? {}) as Record<string, unknown>;

      const stepOutputs = await this.fetchStepOutputs(workflowRunId, dependsOn || []);

      // Build context: merged input + step-level input + previous step outputs
      const context = {
        input: {
          ...mergedInput,
          ...stepLevelInput,
        },
        steps: stepOutputs,
      };

      const requestDef = stepDef.request ?? {};
      const { resolved: resolvedRequest, referencedSteps } = this.templateResolver.resolve(requestDef, context);

      const request = this.applyHttpOverrides(
        stepDef.type,
        resolvedRequest,
        job.data.requestOverride,
      );

      const executor = this.executorRegistry.get(stepDef.type);

      const output = await executor.execute({
        request,
        input: context,
      });

      const stepRun = await this.stepRunRepository.findById(stepRunId);
      if (!stepRun) {
        throw new Error(`StepRun not found: ${stepRunId}`);
      }

      const durationMs = stepRun.startedAt ? Date.now() - stepRun.startedAt.getTime() : 0;

      await this.stepRunRepository.update(stepRunId, {
        status: 'SUCCEEDED',
        finishedAt: new Date(),
        output: output as unknown as object,
        input: context.input as unknown as object,
        durationMs,
      });

      this.logger.log(`Step ${stepKey} completed successfully (${durationMs}ms)`);

      await this.checkWorkflowCompletion(workflowRunId, stepRunId);
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'unknown error';
      this.logger.error(`Step ${stepKey} failed: ${errorMessage}`);

      await this.stepRunRepository.update(stepRunId, {
        status: 'FAILED',
        finishedAt: new Date(),
        error: {
          message: errorMessage,
          stack: error instanceof Error ? error.stack : undefined,
        },
      });

      await this.checkWorkflowCompletion(workflowRunId, stepRunId);

      throw error;
    }
  }

  private async fetchStepOutputs(workflowRunId: string, stepKeys: string[]): Promise<Record<string, unknown>> {
    const outputs: Record<string, unknown> = {};

    for (const key of stepKeys) {
      const stepRun = await this.stepRunRepository.findFirst({
        where: { workflowRunId, stepKey: key, status: 'SUCCEEDED' },
        orderBy: { createdAt: 'desc' },
      });

      if (stepRun?.output) {
        outputs[key] = stepRun.output;
      }
    }

    return outputs;
  }

  private async checkWorkflowCompletion(
    workflowRunId: string,
    completedStepRunId: string,
  ): Promise<void> {
    const siblingSteps = await this.stepRunRepository.findManyByWorkflowRun(workflowRunId);

    const allDone = siblingSteps.every((s) => s.status === 'SUCCEEDED' || s.status === 'FAILED');
    if (!allDone) return;

    const hasFailure = siblingSteps.some((s) => s.status === 'FAILED');
    const finalStatus = hasFailure ? 'FAILED' : 'SUCCEEDED';

    await this.workflowRunRepository.update(workflowRunId, {
      status: finalStatus,
      finishedAt: new Date(),
    });

    this.logger.log(`WorkflowRun ${workflowRunId} completed with status: ${finalStatus}`);
  }

  private applyHttpOverrides(
    stepType: string,
    baseRequest: unknown,
    requestOverride: StepRunJobPayload['requestOverride'],
  ): unknown {
    if (stepType !== 'http' || !requestOverride) return baseRequest;

    if (typeof baseRequest !== 'object' || baseRequest === null) return baseRequest;

    const base = baseRequest as Record<string, unknown>;
    const merged: Record<string, unknown> = { ...base };

    if (requestOverride.query && typeof merged.query === 'object' && merged.query !== null) {
      merged.query = {
        ...(merged.query as Record<string, unknown>),
        ...requestOverride.query,
      };
    } else if (requestOverride.query) {
      merged.query = requestOverride.query;
    }

    if (requestOverride.body !== undefined) {
      merged.body = requestOverride.body;
    }

    return merged;
  }
}
