import type { Event } from '@prisma/client';

type EventOverrides = Partial<Event>;

export const createEventFixture = (overrides: EventOverrides = {}): Event => {
  const now = new Date('2026-01-23T10:00:00.000Z');

  return {
    id: 'ev_1',
    triggerId: 'tr_1',
    type: 'WEBHOOK',
    externalId: 'ext_1',
    payload: {},
    receivedAt: now,
    createdAt: now,
    ...overrides,
  };
};

export const createEventListFixture = (count = 3): Event[] =>
  Array.from({ length: count }, (_, i) =>
    createEventFixture({
      id: `ev_${i + 1}`,
      externalId: `ext_${i + 1}`,
    }),
  );
