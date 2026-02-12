import { Test, TestingModule } from '@nestjs/testing';
import { EventController } from './event.controller';
import { EventService } from './event.service';
import { createEventFixture } from 'test/event/event.fixtures';

describe('EventController', () => {
  let controller: EventController;
  let service: EventService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [EventController],
      providers: [
        {
          provide: EventService,
          useValue: {
            list: jest.fn(),
            get: jest.fn(),
          },
        },
      ],
    }).compile();

    controller = module.get<EventController>(EventController);
    service = module.get<EventService>(EventService);
  });

  it('list() calls EventService.list()', async () => {
    const list = [createEventFixture({ id: 'ev_1' })];
    const listSpy = jest.spyOn(service, 'list').mockResolvedValue({
      items: list,
      pagination: {
        page: 1,
        pageSize: 25,
        total: 1,
        totalPages: 1,
        hasNext: false,
        hasPrev: false,
      },
    });

    await expect(
      controller.list('wf_1', 'tr_1', { page: 1, pageSize: 25 }),
    ).resolves.toEqual({
      items: list,
      pagination: {
        page: 1,
        pageSize: 25,
        total: 1,
        totalPages: 1,
        hasNext: false,
        hasPrev: false,
      },
    });
    expect(listSpy).toHaveBeenCalledWith({
      workflowId: 'wf_1',
      triggerId: 'tr_1',
      page: 1,
      pageSize: 25,
    });
  });

  it('get() calls EventService.get()', async () => {
    const event = createEventFixture({ id: 'ev_1', triggerId: 'tr_1' });
    const getSpy = jest.spyOn(service, 'get').mockResolvedValue(event);

    await expect(controller.get('wf_1', 'tr_1', 'ev_1')).resolves.toBe(event);
    expect(getSpy).toHaveBeenCalledWith('wf_1', 'tr_1', 'ev_1');
  });
});
