import {
  CallHandler,
  ExecutionContext,
  Injectable,
  NestInterceptor,
} from '@nestjs/common';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import type { FastifyReply, FastifyRequest } from 'fastify';

type Envelope = {
  ok: true;
  statusCode: number;
  path: string;
  timestamp: string;
  data: unknown;
  meta: { requestId?: string };
};

function isEnvelope(
  v: unknown,
): v is { ok: unknown; data?: unknown; error?: unknown } {
  return typeof v === 'object' && v !== null && 'ok' in v;
}

@Injectable()
export class ResponseInterceptor implements NestInterceptor {
  intercept(context: ExecutionContext, next: CallHandler): Observable<unknown> {
    const http = context.switchToHttp();
    const req = http.getRequest<FastifyRequest>();
    const res = http.getResponse<FastifyReply>();

    return next.handle().pipe(
      map((data: unknown) => {
        // Don’t double wrap
        if (isEnvelope(data) && ('data' in data || 'error' in data)) {
          return data;
        }

        const body: Envelope = {
          ok: true,
          statusCode: res.statusCode,
          path: req.url,
          timestamp: new Date().toISOString(),
          data,
          meta: { requestId: req.id },
        };

        return body;
      }),
    );
  }
}
