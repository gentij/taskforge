import { Module } from '@nestjs/common';
import { HttpExecutor } from './http-executor';

@Module({
  providers: [HttpExecutor],
  exports: [HttpExecutor],
})
export class HttpExecutorModule {}
