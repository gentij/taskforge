import { Body, Controller, Get, Post } from '@nestjs/common';
import { AppService } from './app.service';

import { createZodDto } from 'nestjs-zod';
import { z } from 'zod';

const CredentialsSchema = z.object({
  username: z.string().min(1),
  password: z.string().min(1),
});

// class is required for using DTO as a type
class CredentialsDto extends createZodDto(CredentialsSchema) {}

@Controller()
export class AppController {
  constructor(private readonly appService: AppService) {}

  @Get()
  getHello(): string {
    return this.appService.getHello();
  }

  @Post()
  validationPipeTest(@Body() _credentials: CredentialsDto): string {
    console.log(_credentials);
    return 'working';
  }
}
