import { Module } from '@nestjs/common';
import { AppLifecycleService } from './app-lifecycle.service';

@Module({
  providers: [AppLifecycleService],
})
export class CoreModule {}
