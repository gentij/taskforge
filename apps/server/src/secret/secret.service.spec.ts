import { Test } from '@nestjs/testing';
import { SecretService } from './secret.service';
import { SecretRepository } from '@taskforge/db-access';
import {
  createSecretRepositoryMock,
  type SecretRepositoryMock,
} from 'test/secret/secret.repository.mock';
import {
  createSecretFixture,
  createSecretListFixture,
} from 'test/secret/secret.fixtures';
import { AppError } from 'src/common/http/errors/app-error';

describe('SecretService', () => {
  let service: SecretService;
  let repo: SecretRepositoryMock;

  beforeEach(async () => {
    repo = createSecretRepositoryMock();

    const moduleRef = await Test.createTestingModule({
      providers: [SecretService, { provide: SecretRepository, useValue: repo }],
    }).compile();

    service = moduleRef.get(SecretService);
  });

  it('create() creates a secret', async () => {
    const created = createSecretFixture({ name: 'API_KEY' });
    repo.create.mockResolvedValue(created);

    await expect(
      service.create({ name: 'API_KEY', value: 'secret-value' }),
    ).resolves.toBe(created);

    expect(repo.create).toHaveBeenCalledWith({
      name: 'API_KEY',
      value: 'secret-value',
      description: undefined,
    });
  });

  it('list() returns secrets', async () => {
    const list = createSecretListFixture(2);
    repo.findMany.mockResolvedValue(list);

    await expect(service.list()).resolves.toBe(list);
    expect(repo.findMany).toHaveBeenCalledTimes(1);
  });

  it('get() returns secret when found', async () => {
    const secret = createSecretFixture({ id: 'sec_1' });
    repo.findById.mockResolvedValue(secret);

    await expect(service.get('sec_1')).resolves.toBe(secret);
    expect(repo.findById).toHaveBeenCalledWith('sec_1');
  });

  it('get() throws notFound when secret missing', async () => {
    repo.findById.mockResolvedValue(null);

    await expect(service.get('missing')).rejects.toBeInstanceOf(AppError);
  });

  it('update() updates secret after existence check', async () => {
    const secret = createSecretFixture({ id: 'sec_1' });
    const updated = createSecretFixture({ id: 'sec_1', name: 'UPDATED' });

    repo.findById.mockResolvedValue(secret);
    repo.update.mockResolvedValue(updated);

    await expect(service.update('sec_1', { name: 'UPDATED' })).resolves.toBe(
      updated,
    );

    expect(repo.update).toHaveBeenCalledWith('sec_1', { name: 'UPDATED' });
  });

  it('delete() deletes secret after existence check', async () => {
    const secret = createSecretFixture({ id: 'sec_1' });
    repo.findById.mockResolvedValue(secret);
    repo.delete.mockResolvedValue(secret);

    await expect(service.delete('sec_1')).resolves.toBe(secret);
    expect(repo.delete).toHaveBeenCalledWith('sec_1');
  });
});
