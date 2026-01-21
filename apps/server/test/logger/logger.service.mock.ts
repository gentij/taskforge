import { PinoLogger } from 'nestjs-pino';

export type PinoLoggerMock = jest.Mocked<
  Pick<
    PinoLogger,
    'setContext' | 'info' | 'warn' | 'error' | 'debug' | 'trace' | 'fatal'
  >
>;

export const createPinoLoggerMock = (): PinoLoggerMock => ({
  setContext: jest.fn(),
  info: jest.fn(),
  warn: jest.fn(),
  error: jest.fn(),
  debug: jest.fn(),
  trace: jest.fn(),
  fatal: jest.fn(),
});
