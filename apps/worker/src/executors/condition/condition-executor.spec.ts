/* eslint-disable @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-member-access */

import { ConditionExecutor } from './condition-executor';

describe('ConditionExecutor', () => {
  it('passes when expr is truthy', async () => {
    const exec = new ConditionExecutor();

    const out = await exec.execute({
      request: {
        expr: 'length(source.items) > `1`',
        source: { items: [1, 2, 3] },
      },
      input: { input: {}, steps: {} },
    });

    expect(out.statusCode).toBe(200);
    expect(out.body).toEqual({ passed: true, value: true });
  });

  it('fails when expr is falsy and assert is true (default)', async () => {
    const exec = new ConditionExecutor();

    await expect(
      exec.execute({
        request: {
          expr: 'length(source.items) > `10`',
          source: { items: [1, 2] },
          message: 'not enough items',
        },
        input: { input: {}, steps: {} },
      }),
    ).rejects.toThrow('Condition failed: not enough items');
  });

  it('returns passed=false when assert is false', async () => {
    const exec = new ConditionExecutor();

    const out = await exec.execute({
      request: {
        expr: 'length(source.items) > `10`',
        assert: false,
        source: { items: [1, 2] },
      },
      input: { input: {}, steps: {} },
    });

    expect(out.body).toEqual({ passed: false, value: false });
  });

  it('treats malformed expr as falsy (and then assert fails)', async () => {
    const exec = new ConditionExecutor();

    await expect(
      exec.execute({
        request: {
          expr: 'length(',
        },
        input: { input: {}, steps: {} },
      }),
    ).rejects.toThrow('Condition failed');
  });
});
