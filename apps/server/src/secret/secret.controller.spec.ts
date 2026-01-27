import { Test, TestingModule } from '@nestjs/testing';
import { SecretController } from './secret.controller';
import { SecretService } from './secret.service';
import { createSecretFixture } from 'test/secret/secret.fixtures';

describe('SecretController', () => {
  let controller: SecretController;
  let service: SecretService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [SecretController],
      providers: [
        {
          provide: SecretService,
          useValue: {
            create: jest.fn(),
            list: jest.fn(),
            get: jest.fn(),
            update: jest.fn(),
            delete: jest.fn(),
          },
        },
      ],
    }).compile();

    controller = module.get<SecretController>(SecretController);
    service = module.get<SecretService>(SecretService);
  });

  it('create() calls SecretService.create()', async () => {
    const secret = createSecretFixture({ name: 'API_KEY' });
    const createSpy = jest.spyOn(service, 'create').mockResolvedValue(secret);

    await expect(
      controller.create({ name: 'API_KEY', value: 'secret-value' }),
    ).resolves.toBe(secret);

    expect(createSpy).toHaveBeenCalledWith({
      name: 'API_KEY',
      value: 'secret-value',
      description: undefined,
    });
  });

  it('list() calls SecretService.list()', async () => {
    const list = [createSecretFixture({ id: 'sec_1' })];
    const listSpy = jest.spyOn(service, 'list').mockResolvedValue(list);

    await expect(controller.list()).resolves.toBe(list);
    expect(listSpy).toHaveBeenCalledTimes(1);
  });

  it('get() calls SecretService.get()', async () => {
    const secret = createSecretFixture({ id: 'sec_1' });
    const getSpy = jest.spyOn(service, 'get').mockResolvedValue(secret);

    await expect(controller.get('sec_1')).resolves.toBe(secret);
    expect(getSpy).toHaveBeenCalledWith('sec_1');
  });

  it('update() calls SecretService.update()', async () => {
    const secret = createSecretFixture({ id: 'sec_1', name: 'UPDATED' });
    const updateSpy = jest.spyOn(service, 'update').mockResolvedValue(secret);

    await expect(controller.update('sec_1', { name: 'UPDATED' })).resolves.toBe(
      secret,
    );

    expect(updateSpy).toHaveBeenCalledWith('sec_1', {
      name: 'UPDATED',
      value: undefined,
      description: undefined,
    });
  });

  it('delete() calls SecretService.delete()', async () => {
    const secret = createSecretFixture({ id: 'sec_1' });
    const deleteSpy = jest.spyOn(service, 'delete').mockResolvedValue(secret);

    await expect(controller.delete('sec_1')).resolves.toBe(secret);
    expect(deleteSpy).toHaveBeenCalledWith('sec_1');
  });
});
