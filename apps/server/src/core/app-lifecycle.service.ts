import { Injectable, Logger, OnApplicationShutdown } from '@nestjs/common';

@Injectable()
export class AppLifecycleService implements OnApplicationShutdown {
  private readonly logger = new Logger(AppLifecycleService.name);

  onApplicationShutdown(signal?: string) {
    this.logger.log(`Shutting down (signal: ${signal ?? 'unknown'})`);
    // Later:
    // - await prisma.$disconnect()
    // - await redis.quit()
    // - await worker.stop()
  }
}
