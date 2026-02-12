import { Injectable } from '@nestjs/common';
import { ThrottlerGuard } from '@nestjs/throttler';
type TrackerReq = {
  ip?: string;
  headers?: Record<string, unknown>;
};

@Injectable()
export class ThrottlerBehindProxyGuard extends ThrottlerGuard {
  // eslint-disable-next-line @typescript-eslint/require-await
  protected async getTracker(req: TrackerReq): Promise<string> {
    const ip = typeof req.ip === 'string' ? req.ip : undefined;

    const rawXff = req.headers?.['x-forwarded-for'];
    let forwarded: string | undefined;
    if (Array.isArray(rawXff)) {
      forwarded = typeof rawXff[0] === 'string' ? rawXff[0] : undefined;
    } else if (typeof rawXff === 'string') {
      forwarded = rawXff;
    }

    const firstForwarded = forwarded
      ? forwarded.split(',')[0].trim()
      : undefined;

    return ip ?? firstForwarded ?? 'unknown';
  }
}
