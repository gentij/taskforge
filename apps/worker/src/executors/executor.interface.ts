export interface ExecutorOutput {
  statusCode: number;
  headers: Record<string, string> | undefined;
  body: unknown;
}

export interface StepExecutor {
  readonly stepType: string;
  execute(input: unknown): Promise<ExecutorOutput>;
}
