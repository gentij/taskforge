-- CreateTable
CREATE TABLE "Event" (
    "id" TEXT NOT NULL,
    "triggerId" TEXT NOT NULL,
    "type" TEXT,
    "externalId" TEXT,
    "payload" JSONB NOT NULL DEFAULT '{}',
    "receivedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "Event_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "Event_triggerId_idx" ON "Event"("triggerId");

-- CreateIndex
CREATE INDEX "Event_receivedAt_idx" ON "Event"("receivedAt");

-- CreateIndex
CREATE UNIQUE INDEX "Event_triggerId_externalId_key" ON "Event"("triggerId", "externalId");

-- AddForeignKey
ALTER TABLE "Event" ADD CONSTRAINT "Event_triggerId_fkey" FOREIGN KEY ("triggerId") REFERENCES "Trigger"("id") ON DELETE CASCADE ON UPDATE CASCADE;
