import type { Config } from 'jest';

const config: Config = {
  rootDir: '.',

  testEnvironment: 'node',

  testRegex: '.*\\.spec\\.ts$',

  transform: {
    '^.+\\.(t|j)s$': ['ts-jest', { tsconfig: '<rootDir>/tsconfig.json' }],
  },

  moduleFileExtensions: ['ts', 'js', 'json'],

  moduleNameMapper: {
    '^src/(.*)$': '<rootDir>/src/$1',
    '^test/(.*)$': '<rootDir>/test/$1',
  },

  clearMocks: true,
};

export default config;
