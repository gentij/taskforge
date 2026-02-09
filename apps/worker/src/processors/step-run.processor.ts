import { Processor, WorkerHost } from '@nestjs/bullmq';
import { Logger } from '@nestjs/common';
import { Job } from 'bullmq';
import {
  StepRunRepository,
  WorkflowRunRepository,
  WorkflowVersionRepository,
  SecretRepository,
} from '@taskforge/db-access';
import { StepRunJobPayload } from '@taskforge/contracts';
import { ExecutorRegistry } from '../executors/executor-registry';
import { TemplateResolver } from '../utils/template-resolver';
import { wrapForDb } from '../utils/persisted-json';
import { redactSecrets } from '../utils/redact';
import { CryptoService } from '../crypto/crypto.service';

interface StepDefinition {
  key: string;
  type: string;
  dependsOn?: string[];
  input?: Record<string, unknown>;
  request?: Record<string, unknown>;
  outputPolicy?: {
    truncate?: boolean;
    maxBytes?: number;
  };
}

@Processor('step-runs')
export class StepRunProcessor extends WorkerHost {
  private readonly logger = new Logger(StepRunProcessor.name);
  private readonly templateResolver = new TemplateResolver();

  constructor(
    private readonly stepRunRepository: StepRunRepository,
    private readonly workflowRunRepository: WorkflowRunRepository,
    private readonly workflowVersionRepository: WorkflowVersionRepository,
    private readonly secretRepository: SecretRepository,
    private readonly crypto: CryptoService,
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

      const secretNames = this.findSecretRefs(stepDef.request ?? {});
      const secrets = await this.secretRepository.findManyByNames(secretNames);
      const secretMap: Record<string, string> = {};
      for (const s of secrets) secretMap[s.name] = this.crypto.decryptSecret(s.value);
      const secretValues = Object.values(secretMap);

      // Build context: merged input + step-level input + previous step outputs
      const context = {
        input: {
          ...mergedInput,
          ...stepLevelInput,
        },
        steps: stepOutputs,
        secret: secretMap,
      };

      const requestDef = stepDef.request ?? {};
      const { resolved: resolvedRequest, referencedSteps } = this.templateResolver.resolve(requestDef, context);

      const request = this.applyHttpOverrides(
        stepDef.type,
        resolvedRequest,
        job.data.requestOverride,
      );

      const executor = this.executorRegistry.get(stepDef.type);

      const rawOutput = await executor.execute({
        request,
        input: context,
      });

      const outputUnwrapped = unwrapExecutorOutput(rawOutput);

      const stepRun = await this.stepRunRepository.findById(stepRunId);
      if (!stepRun) {
        throw new Error(`StepRun not found: ${stepRunId}`);
      }

      const durationMs = stepRun.startedAt ? Date.now() - stepRun.startedAt.getTime() : 0;

      const outputPolicy = stepDef.outputPolicy ?? {};
      const redactedOutput = redactSecrets(outputUnwrapped, { secretValues });
      const outputEnvelope = wrapForDb(redactedOutput, {
        maxBytes: outputPolicy.maxBytes ?? 256 * 1024,
        truncate: outputPolicy.truncate ?? true,
        reason: 'step_output',
      });

      const redactedInput = redactSecrets(context.input, { secretValues });
      const inputEnvelope = wrapForDb(redactedInput, {
        maxBytes: 256 * 1024,
        truncate: true,
        reason: 'step_input',
      });

      await this.stepRunRepository.update(stepRunId, {
        status: 'SUCCEEDED',
        finishedAt: new Date(),
        output: outputEnvelope as unknown as object,
        input: inputEnvelope as unknown as object,
        durationMs,
      });

      this.logger.log(`Step ${stepKey} completed successfully (${durationMs}ms)`);

      await this.checkWorkflowCompletion(workflowRunId, stepRunId);
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'unknown error';
      this.logger.error(`Step ${stepKey} failed: ${errorMessage}`);

      // Best-effort: if we already loaded secrets earlier in the run, they are not in scope here.
      // We still redact common sensitive keys.

      const redactedError = redactSecrets(
        {
          message: errorMessage,
          stack: error instanceof Error ? error.stack : undefined,
        },
        { secretValues: [] },
      );
      const errorEnvelope = wrapForDb(redactedError, {
        maxBytes: 64 * 1024,
        truncate: true,
        reason: 'step_error',
      });

      await this.stepRunRepository.update(stepRunId, {
        status: 'FAILED',
        finishedAt: new Date(),
        lastErrorAt: new Date(),
        error: errorEnvelope as unknown as object,
      });

      await this.checkWorkflowCompletion(workflowRunId, stepRunId);

      throw error;
    }
  }

  private findSecretRefs(value: unknown): string[] {
    const found = new Set<string>();
    const pattern = /\{\{\s*secret\.([a-zA-Z0-9_-]+)\s*\}\}/g;

    const walk = (node: unknown) => {
      if (typeof node === 'string') {
        pattern.lastIndex = 0;
        let m: RegExpExecArray | null;
        while ((m = pattern.exec(node)) !== null) {
          if (m[1]) found.add(m[1]);
        }
        return;
      }
      if (Array.isArray(node)) {
        for (const x of node) walk(x);
        return;
      }
      if (node && typeof node === 'object') {
        for (const v of Object.values(node as Record<string, unknown>)) walk(v);
      }
    };

    walk(value);
    return Array.from(found);
  }

  private async fetchStepOutputs(workflowRunId: string, stepKeys: string[]): Promise<Record<string, unknown>> {
    const outputs: Record<string, unknown> = {};

    for (const key of stepKeys) {
      const stepRun = await this.stepRunRepository.findFirst({
        where: { workflowRunId, stepKey: key, status: 'SUCCEEDED' },
        orderBy: { createdAt: 'desc' },
      });

      if (stepRun?.output) {
        const o = stepRun.output as unknown as { data?: unknown };
        outputs[key] = o && typeof o === 'object' && 'data' in o ? (o as any).data : stepRun.output;
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

function unwrapExecutorOutput(value: unknown): unknown {
  if (!value || typeof value !== 'object') return value;
  // Unwrap HttpExecutor body helper wrapper.
  if ('_taskforgeHttp' in (value as any) && 'data' in (value as any)) {
    return (value as any).data;
  }
  return value;
}
