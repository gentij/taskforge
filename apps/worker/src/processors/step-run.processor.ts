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

@Processor('step-runs')
export class StepRunProcessor extends WorkerHost {
  private readonly logger = new Logger(StepRunProcessor.name);

  constructor(
    private readonly stepRunRepository: StepRunRepository,
    private readonly workflowRunRepository: WorkflowRunRepository,
    private readonly workflowVersionRepository: WorkflowVersionRepository,
    private readonly executorRegistry: ExecutorRegistry,
  ) {
    super();
  }

  async process(job: Job<StepRunJobPayload>): Promise<void> {
    const { workflowRunId, stepRunId, stepKey, workflowVersionId, input } = job.data;

    this.logger.log(`Processing step ${stepKey} for run ${workflowRunId} (jobId: ${job.id})`);

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
        steps: Array<{ key: string; type: string; request: unknown }>;
      };
      const stepDef = definition.steps.find((s) => s.key === stepKey);

      if (!stepDef) {
        throw new Error(`Step definition not found for key: ${stepKey}`);
      }

      const executor = this.executorRegistry.get(stepDef.type);

      const output = await executor.execute({
        request: stepDef.request,
        input,
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
        durationMs,
      });

      this.logger.log(`Step ${stepKey} completed successfully (${durationMs}ms)`);

      await this.checkWorkflowCompletion(workflowRunId, stepRunId);
    } catch (error) {
      this.logger.error(
        `Step ${stepKey} failed: ${error instanceof Error ? error.message : 'unknown error'}`,
      );

      await this.stepRunRepository.update(stepRunId, {
        status: 'FAILED',
        finishedAt: new Date(),
        error: {
          message: error instanceof Error ? error.message : 'unknown error',
          stack: error instanceof Error ? error.stack : undefined,
        },
      });

      await this.checkWorkflowCompletion(workflowRunId, stepRunId);

      throw error;
    }
  }

  private async checkWorkflowCompletion(
    workflowRunId: string,
    completedStepRunId: string,
  ): Promise<void> {
    void completedStepRunId;

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
}
