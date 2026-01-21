export const ErrorDefinitions = {
  COMMON: {
    VALIDATION_ERROR: {
      code: 'VALIDATION_ERROR',
      message: 'Request validation failed',
    },
    INTERNAL_ERROR: {
      code: 'INTERNAL_ERROR',
      message: 'Something went wrong',
    },
    UNAUTHORIZED: {
      code: 'UNAUTHORIZED',
      message: 'Unauthorized',
    },
    FORBIDDEN: {
      code: 'FORBIDDEN',
      message: 'Forbidden',
    },
    RATE_LIMITED: {
      code: 'RATE_LIMITED',
      message: 'Too many requests',
    },
    NOT_FOUND: {
      code: 'NOT_FOUND',
      message: 'Not found',
    },
    BAD_REQUEST: {
      code: 'BAD_REQUEST',
      message: 'Bad request.',
    },
    CONFLICT: { code: 'CONFLICT', message: 'Conflict' },
  },

  AUTH: {
    MISSING_BEARER_TOKEN: {
      code: 'AUTH_MISSING_BEARER_TOKEN',
      message: 'Missing Authorization bearer token',
    },
    INVALID_TOKEN: {
      code: 'AUTH_INVALID_TOKEN',
      message: 'Invalid API token',
    },
    REVOKED_TOKEN: {
      code: 'AUTH_REVOKED_TOKEN',
      message: 'API token has been revoked',
    },
  },

  WORKFLOW: {
    NOT_FOUND: {
      code: 'WORKFLOW_NOT_FOUND',
      message: 'Workflow not found',
    },
    INVALID_STATE: {
      code: 'WORKFLOW_INVALID_STATE',
      message: 'Invalid workflow state',
    },
  },

  STEP: {
    NOT_FOUND: {
      code: 'STEP_NOT_FOUND',
      message: 'Step not found',
    },
  },

  DATABASE: {
    UNIQUE_CONSTRAINT: {
      code: 'UNIQUE_CONSTRAINT',
      message: 'Resource already exists',
    },
  },
} as const;
