import { Injectable, Logger } from '@nestjs/common';
import { StepExecutor } from './executor.interface';
import { HttpExecutor } from './http/http-executor';

@Injectable()
export class ExecutorRegistry {
  private readonly logger = new Logger(ExecutorRegistry.name);
  private readonly executors: Map<string, StepExecutor> = new Map();

  constructor(private readonly httpExecutor: HttpExecutor) {
    this.register(httpExecutor);
  }

  private register(executor: StepExecutor): void {
    this.executors.set(executor.stepType, executor);
    this.logger.log(`Registered executor for step type: ${executor.stepType}`);
  }

  get(stepType: string): StepExecutor {
    const executor = this.executors.get(stepType);

    if (!executor) {
      this.logger.error(`No executor found for step type: ${stepType}`);
      throw new Error(`No executor found for step type: ${stepType}`);
    }

    return executor;
  }

  getRegisteredTypes(): string[] {
    return Array.from(this.executors.keys());
  }
}
