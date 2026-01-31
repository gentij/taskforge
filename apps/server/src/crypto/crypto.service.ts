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

  public generateId(): string {
    const bytes = new Uint8Array(16);
    crypto.getRandomValues(bytes);
    return Array.from(bytes, (b) => b.toString(16).padStart(2, '0')).join('');
  }
}

function timingSafeEqualHex(a: string, b: string): boolean {
  if (a.length !== b.length) return false;
  return timingSafeEqual(Buffer.from(a), Buffer.from(b));
}
