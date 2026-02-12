import { TransformExecutor } from './transform-executor';

describe('TransformExecutor', () => {
  it('evaluates $jmes nodes against input + source', async () => {
    const exec = new TransformExecutor();

    const out = await exec.execute({
      request: {
        source: {
          users: [
            { id: 1, email: 'a@test.com' },
            { id: 2, email: 'b@test.com' },
          ],
        },
        output: {
          usersCount: { $jmes: 'length(source.users)' },
          userEmails: { $jmes: 'source.users[].email' },
        },
      },
      input: {
        input: { apiUrl: 'https://example.com' },
        steps: {},
      },
    });

    expect(out.statusCode).toBe(200);
    expect(out.body).toEqual({
      usersCount: 2,
      userEmails: ['a@test.com', 'b@test.com'],
    });
  });

  it('keeps single-item arrays as arrays', async () => {
    const exec = new TransformExecutor();

    const out = await exec.execute({
      request: {
        source: {
          users: [
            { id: 1, name: 'Alice' },
            { id: 2, name: 'Bob' },
          ],
        },
        output: {
          filteredNames: { $jmes: 'source.users[?id == `1`].name' },
        },
      },
      input: {
        input: {},
        steps: {},
      },
    });

    expect(out.body).toEqual({
      filteredNames: ['Alice'],
    });
  });

  it('exposes steps as unwrapped bodies and stepResponses as full outputs', async () => {
    const exec = new TransformExecutor();

    const out = await exec.execute({
      request: {
        output: {
          userName: { $jmes: 'steps.fetchUser.name' },
          status: { $jmes: 'stepResponses.fetchUser.statusCode' },
        },
      },
      input: {
        input: {},
        steps: {
          fetchUser: { statusCode: 200, body: { id: 1, name: 'Alice' } },
        },
      },
    });

    expect(out.body).toEqual({ userName: 'Alice', status: 200 });
  });
});
