/*
  Warnings:

  - You are about to drop the column `user_one_id` on the `friend` table. All the data in the column will be lost.
  - You are about to drop the column `user_two_id` on the `friend` table. All the data in the column will be lost.
  - You are about to drop the column `user_id` on the `song_of_the_day` table. All the data in the column will be lost.
  - You are about to drop the column `photo_url` on the `user` table. All the data in the column will be lost.
  - You are about to drop the column `user_id` on the `user_statistic` table. All the data in the column will be lost.
  - A unique constraint covering the columns `[userId,setAt]` on the table `song_of_the_day` will be added. If there are existing duplicate values, this will fail.
  - Added the required column `userOneId` to the `friend` table without a default value. This is not possible if the table is not empty.
  - Added the required column `userTwoId` to the `friend` table without a default value. This is not possible if the table is not empty.
  - Added the required column `userId` to the `song_of_the_day` table without a default value. This is not possible if the table is not empty.
  - Added the required column `userId` to the `user_statistic` table without a default value. This is not possible if the table is not empty.

*/
-- DropForeignKey
ALTER TABLE "friend" DROP CONSTRAINT "friend_user_one_id_fkey";

-- DropForeignKey
ALTER TABLE "friend" DROP CONSTRAINT "friend_user_two_id_fkey";

-- DropForeignKey
ALTER TABLE "song_of_the_day" DROP CONSTRAINT "song_of_the_day_user_id_fkey";

-- DropForeignKey
ALTER TABLE "user_statistic" DROP CONSTRAINT "user_statistic_user_id_fkey";

-- DropIndex
DROP INDEX "song_of_the_day_user_id_setAt_key";

-- AlterTable
ALTER TABLE "friend" DROP COLUMN "user_one_id",
DROP COLUMN "user_two_id",
ADD COLUMN     "userOneId" TEXT NOT NULL,
ADD COLUMN     "userTwoId" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "song_of_the_day" DROP COLUMN "user_id",
ADD COLUMN     "userId" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "user" DROP COLUMN "photo_url",
ADD COLUMN     "photoUrl" TEXT;

-- AlterTable
ALTER TABLE "user_statistic" DROP COLUMN "user_id",
ADD COLUMN     "userId" TEXT NOT NULL;

-- CreateIndex
CREATE UNIQUE INDEX "song_of_the_day_userId_setAt_key" ON "song_of_the_day"("userId", "setAt");

-- AddForeignKey
ALTER TABLE "friend" ADD CONSTRAINT "friend_userOneId_fkey" FOREIGN KEY ("userOneId") REFERENCES "user"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "friend" ADD CONSTRAINT "friend_userTwoId_fkey" FOREIGN KEY ("userTwoId") REFERENCES "user"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "song_of_the_day" ADD CONSTRAINT "song_of_the_day_userId_fkey" FOREIGN KEY ("userId") REFERENCES "user"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "user_statistic" ADD CONSTRAINT "user_statistic_userId_fkey" FOREIGN KEY ("userId") REFERENCES "user"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
