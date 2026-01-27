-- CreateEnum
CREATE TYPE "StepRunStatus" AS ENUM ('QUEUED', 'RUNNING', 'SUCCEEDED', 'FAILED');

-- CreateTable
CREATE TABLE "StepRun" (
    "id" TEXT NOT NULL,
    "workflowRunId" TEXT NOT NULL,
    "stepKey" TEXT NOT NULL,
    "status" "StepRunStatus" NOT NULL DEFAULT 'QUEUED',
    "attempt" INTEGER NOT NULL DEFAULT 0,
    "input" JSONB NOT NULL DEFAULT '{}',
    "output" JSONB,
    "error" JSONB,
    "logs" JSONB,
    "lastErrorAt" TIMESTAMP(3),
    "durationMs" INTEGER,
    "startedAt" TIMESTAMP(3),
    "finishedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "StepRun_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "StepRun_workflowRunId_idx" ON "StepRun"("workflowRunId");

-- CreateIndex
CREATE INDEX "StepRun_status_idx" ON "StepRun"("status");

-- CreateIndex
CREATE INDEX "StepRun_createdAt_idx" ON "StepRun"("createdAt");

-- CreateIndex
CREATE UNIQUE INDEX "StepRun_workflowRunId_stepKey_key" ON "StepRun"("workflowRunId", "stepKey");

-- AddForeignKey
ALTER TABLE "StepRun" ADD CONSTRAINT "StepRun_workflowRunId_fkey" FOREIGN KEY ("workflowRunId") REFERENCES "WorkflowRun"("id") ON DELETE CASCADE ON UPDATE CASCADE;
