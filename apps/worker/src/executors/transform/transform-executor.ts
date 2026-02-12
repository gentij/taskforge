import { Injectable, Logger } from '@nestjs/common';
import * as jmespath from 'jmespath';
import type { ExecutorOutput, StepExecutor } from '../executor.interface';
import { TransformExecutorInputSchema } from './transform.types';

type RootContext = {
  input: unknown;
  source: Record<string, unknown>;
  steps: Record<string, unknown>;
  stepResponses: Record<string, unknown>;
};

@Injectable()
export class TransformExecutor implements StepExecutor {
  readonly stepType = 'transform';
  private readonly logger = new Logger(TransformExecutor.name);

  async execute(input: unknown): Promise<ExecutorOutput> {
    const validated = TransformExecutorInputSchema.parse(input);
    const request = validated.request;
    const ctx = validated.input as any;

    const source = request.source ?? {};
    const stepResponses = (ctx?.steps ?? {}) as Record<string, unknown>;
    const steps = this.unwrapBodies(stepResponses);

    const root: RootContext = {
      input: ctx?.input ?? {},
      source,
      steps,
      stepResponses,
    };

    const body = this.evaluateTemplate(request.output, root);

    return {
      statusCode: 200,
      headers: undefined,
      body,
    };
  }

  private evaluateTemplate(node: unknown, root: RootContext): unknown {
    if (Array.isArray(node)) {
      return node.map((item) => this.evaluateTemplate(item, root));
    }

    if (node && typeof node === 'object') {
      const obj = node as Record<string, unknown>;
      const keys = Object.keys(obj);

      if (keys.length === 1 && keys[0] === '$jmes') {
        const expr = obj.$jmes;
        if (typeof expr !== 'string' || expr.trim().length === 0) {
          throw new Error('$jmes expression must be a non-empty string');
        }

        try {
          return jmespath.search(root as any, expr);
        } catch (err) {
          const message = err instanceof Error ? err.message : String(err);
          this.logger.error(`JMESPath evaluation failed: ${message}`);
          throw new Error(`JMESPath evaluation failed: ${message}`);
        }
      }

      const out: Record<string, unknown> = {};
      for (const [k, v] of Object.entries(obj)) {
        out[k] = this.evaluateTemplate(v, root);
      }
      return out;
    }

    return node;
  }

  private unwrapBodies(stepResponses: Record<string, unknown>): Record<string, unknown> {
    const out: Record<string, unknown> = {};

    for (const [key, value] of Object.entries(stepResponses)) {
      if (value && typeof value === 'object' && 'body' in (value as any)) {
        out[key] = (value as any).body;
      } else {
        out[key] = value;
      }
    }

    return out;
  }
}
