import { Module } from '@nestjs/common';
import { TransformExecutor } from './transform-executor';

@Module({
  providers: [TransformExecutor],
  exports: [TransformExecutor],
})
export class TransformExecutorModule {}
