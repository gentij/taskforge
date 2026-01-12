import { applyDecorators, Type } from '@nestjs/common';
import { ApiExtraModels, ApiOkResponse, getSchemaPath } from '@nestjs/swagger';
import { ApiResponseDto } from '../dto/api-response.dto';
import { ApiErrorResponseDto } from '../dto/api-error.dto';
import {
  EnvelopeError,
  errorDecorators,
  errorExamples,
} from './api-envelope.lib';

export const ApiEnvelope = <TModel extends Type<unknown>>(
  model: TModel,
  opts?: { description?: string; errors?: EnvelopeError[] },
): MethodDecorator & ClassDecorator => {
  const errors = opts?.errors ?? [400, 500];

  const errorDecs: Array<MethodDecorator> = errors.map((code) =>
    errorDecorators[code]({
      content: {
        'application/json': {
          schema: { $ref: getSchemaPath(ApiErrorResponseDto) },
          examples: {
            default: {
              summary: `${code} example`,
              value: errorExamples[code],
            },
          },
        },
      },
    }),
  );

  return applyDecorators(
    ApiExtraModels(ApiResponseDto, ApiErrorResponseDto, model),

    ApiOkResponse({
      description: opts?.description,
      schema: {
        allOf: [
          { $ref: getSchemaPath(ApiResponseDto) },
          { properties: { data: { $ref: getSchemaPath(model) } } },
        ],
      },
    }),

    ...errorDecs,
  );
};
