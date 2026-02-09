import { redactSecrets } from './redact';

describe('redactSecrets', () => {
  it('redacts values by sensitive keys', () => {
    const out = redactSecrets(
      { webhookUrl: 'x', Authorization: 'y', nested: { apiKey: 'z' } },
      { secretValues: [] },
    );
    expect(out).toEqual({
      webhookUrl: '[REDACTED]',
      Authorization: '[REDACTED]',
      nested: { apiKey: '[REDACTED]' },
    });
  });

  it('redacts secret literals inside strings', () => {
    const out = redactSecrets(
      { msg: 'hello SECRET world', arr: ['SECRET'] },
      { secretValues: ['SECRET'] },
    );
    expect(out).toEqual({ msg: 'hello [REDACTED] world', arr: ['[REDACTED]'] });
  });
});
