-- AlterTable
ALTER TABLE "StepRun" ADD COLUMN     "requestOverride" JSONB;

-- AlterTable
ALTER TABLE "WorkflowRun" ADD COLUMN     "overrides" JSONB;
