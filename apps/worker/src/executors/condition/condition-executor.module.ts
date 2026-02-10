import { Module } from '@nestjs/common';
import { ConditionExecutor } from './condition-executor';

@Module({
  providers: [ConditionExecutor],
  exports: [ConditionExecutor],
})
export class ConditionExecutorModule {}
