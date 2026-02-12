const DEFAULT_SENSITIVE_KEY_PATTERN = /(webhook|token|secret|password|api[-_]?key|authorization)/i;

export function redactSecrets(value: unknown, opts: { secretValues: string[] }): unknown {
  const { secretValues } = opts;
  const needles = secretValues
    .filter((v) => typeof v === 'string' && v.length > 0)
    .sort((a, b) => b.length - a.length);

  return redactNode(value, needles);
}

function redactNode(node: unknown, secretValues: string[]): unknown {
  if (node === null || node === undefined) return node;

  if (typeof node === 'string') {
    return redactString(node, secretValues);
  }

  if (typeof node === 'number' || typeof node === 'boolean') return node;

  if (Array.isArray(node)) {
    return node.map((x) => redactNode(x, secretValues));
  }

  if (typeof node === 'object') {
    const obj = node as Record<string, unknown>;
    const out: Record<string, unknown> = {};
    for (const [k, v] of Object.entries(obj)) {
      if (DEFAULT_SENSITIVE_KEY_PATTERN.test(k)) {
        out[k] = '[REDACTED]';
      } else {
        out[k] = redactNode(v, secretValues);
      }
    }
    return out;
  }

  return node;
}

function redactString(input: string, secretValues: string[]): string {
  let out = input;
  for (const secret of secretValues) {
    if (!secret) continue;
    if (out.includes(secret)) {
      out = out.split(secret).join('[REDACTED]');
    }
  }
  return out;
}
