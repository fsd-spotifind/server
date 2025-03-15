/*
  Warnings:

  - A unique constraint covering the columns `[user_id,setAt]` on the table `song_of_the_day` will be added. If there are existing duplicate values, this will fail.

*/
-- AlterTable
ALTER TABLE "song_of_the_day" ALTER COLUMN "setAt" DROP DEFAULT,
ALTER COLUMN "setAt" SET DATA TYPE TEXT;

-- CreateIndex
CREATE UNIQUE INDEX "song_of_the_day_user_id_setAt_key" ON "song_of_the_day"("user_id", "setAt");
