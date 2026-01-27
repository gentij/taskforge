import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Patch,
  Post,
} from '@nestjs/common';
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger';
import { ApiEnvelope } from 'src/common/swagger/envelope/api-envelope.decorator';
import { SecretService } from './secret.service';
import {
  CreateSecretReqDto,
  SecretResDto,
  UpdateSecretReqDto,
} from './dto/secret.dto';

@ApiTags('Secrets')
@ApiBearerAuth('bearer')
@Controller('secrets')
export class SecretController {
  constructor(private readonly service: SecretService) {}

  @ApiEnvelope(SecretResDto, {
    description: 'Create secret',
    errors: [401, 409, 500],
  })
  @Post()
  create(@Body() body: CreateSecretReqDto) {
    return this.service.create({
      name: body.name,
      value: body.value,
      description: body.description,
    });
  }

  @ApiEnvelope(SecretResDto, {
    description: 'List secrets',
    isArray: true,
    errors: [401, 500],
  })
  @Get()
  list() {
    return this.service.list();
  }

  @ApiEnvelope(SecretResDto, {
    description: 'Get secret',
    errors: [401, 404, 500],
  })
  @Get(':id')
  get(@Param('id') id: string) {
    return this.service.get(id);
  }

  @ApiEnvelope(SecretResDto, {
    description: 'Update secret',
    errors: [401, 404, 500],
  })
  @Patch(':id')
  update(@Param('id') id: string, @Body() body: UpdateSecretReqDto) {
    return this.service.update(id, {
      name: body.name,
      value: body.value,
      description: body.description,
    });
  }

  @ApiEnvelope(SecretResDto, {
    description: 'Delete secret',
    errors: [401, 404, 500],
  })
  @Delete(':id')
  delete(@Param('id') id: string) {
    return this.service.delete(id);
  }
}
