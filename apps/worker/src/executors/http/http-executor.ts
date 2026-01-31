import { Injectable, Logger } from '@nestjs/common';
import { StepExecutor, ExecutorOutput } from '../executor.interface';
import { HttpExecutorInput, HttpRequestSpec } from './http.types';

@Injectable()
export class HttpExecutor implements StepExecutor {
  readonly stepType = 'http';
  private readonly logger = new Logger(HttpExecutor.name);

  async execute(input: unknown): Promise<ExecutorOutput> {
    const validated = input as HttpExecutorInput;
    const { request } = validated;

    this.logger.debug(`Executing HTTP ${request.method} ${request.url}`);

    const startTime = Date.now();

    try {
      const response = await this.makeRequest(request);
      const durationMs = Date.now() - startTime;

      this.logger.log(
        `HTTP ${request.method} ${request.url} - ${response.statusCode} (${durationMs}ms)`,
      );

      return {
        statusCode: response.statusCode,
        headers: response.headers,
        body: response.body,
      };
    } catch (error) {
      const durationMs = Date.now() - startTime;

      this.logger.error(
        `HTTP ${request.method} ${request.url} failed after ${durationMs}ms: ${error instanceof Error ? error.message : 'unknown error'}`,
      );

      throw error;
    }
  }

  private async makeRequest(request: HttpRequestSpec): Promise<{
    statusCode: number;
    headers: Record<string, string>;
    body: unknown;
  }> {
    const url = new URL(request.url);

    if (request.query) {
      for (const [key, value] of Object.entries(request.query)) {
        url.searchParams.append(key, String(value));
      }
    }

    const headers: Record<string, string> = request.headers ?? {};

    const fetchOptions: RequestInit = {
      method: request.method,
      headers,
      signal: request.timeoutMs ? AbortSignal.timeout(request.timeoutMs) : undefined,
    };

    if (request.body && ['POST', 'PUT', 'PATCH'].includes(request.method)) {
      fetchOptions.body = JSON.stringify(request.body);
      if (!headers['Content-Type']) {
        headers['Content-Type'] = 'application/json';
      }
    }

    const response = await fetch(url.toString(), fetchOptions);

    const responseHeaders: Record<string, string> = {};
    response.headers.forEach((value, key) => {
      responseHeaders[key] = value;
    });

    let body: unknown;
    const contentType = response.headers.get('content-type') ?? '';

    if (contentType.includes('application/json')) {
      body = await response.json();
    } else if (contentType.includes('text/')) {
      body = await response.text();
    } else {
      body = await response.text();
    }

    return {
      statusCode: response.status,
      headers: responseHeaders,
      body,
    };
  }
}
