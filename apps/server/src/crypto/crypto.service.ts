// src/core/crypto/crypto.service.ts
import { Injectable } from '@nestjs/common';
import { randomBytes, createHash, createHmac, timingSafeEqual } from 'crypto';

@Injectable()
export class CryptoService {
  generateApiToken(): string {
    return `tf_${randomBytes(32).toString('hex')}`;
  }

  hashApiToken(token: string): string {
    return createHash('sha256').update(token).digest('hex');
  }

  verifyHmac(
    payload: string | Buffer,
    secret: string,
    signature: string,
  ): boolean {
    const computed = createHmac('sha256', secret).update(payload).digest('hex');

    return timingSafeEqualHex(computed, signature);
  }
}

function timingSafeEqualHex(a: string, b: string): boolean {
  if (a.length !== b.length) return false;
  return timingSafeEqual(Buffer.from(a), Buffer.from(b));
}
