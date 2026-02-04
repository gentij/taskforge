import { TemplateResolver } from './template-resolver';

describe('TemplateResolver', () => {
  const resolver = new TemplateResolver();

  it('interpolates input strings without JSON quotes', () => {
    const { resolved } = resolver.resolve(
      { url: '{{input.apiUrl}}/users', method: 'GET' },
      { input: { apiUrl: 'https://example.com' }, steps: {} },
    );

    expect(resolved).toEqual({ url: 'https://example.com/users', method: 'GET' });
  });

  it('resolves full input expression to a string', () => {
    const { resolved } = resolver.resolve('{{input.apiUrl}}', {
      input: { apiUrl: 'https://example.com' },
      steps: {},
    });

    expect(resolved).toBe('https://example.com');
  });

  it('resolves full step expression to an object (auto-unwrap body)', () => {
    const { resolved } = resolver.resolve('{{steps.fetch.output}}', {
      input: {},
      steps: { fetch: { statusCode: 200, body: { ok: true } } },
    });

    expect(resolved).toEqual({ ok: true });
  });

  it('interpolates primitives inside strings', () => {
    const { resolved } = resolver.resolve('Hello {{input.name}}!', {
      input: { name: 'Alice' },
      steps: {},
    });

    expect(resolved).toBe('Hello Alice!');
  });
});
