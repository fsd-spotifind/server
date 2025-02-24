/*
  Warnings:

  - You are about to drop the `friends` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `users` table. If the table is not empty, all the data it contains will be lost.

*/
-- CreateEnum
CREATE TYPE "StatisticPeriod" AS ENUM ('weekly', 'monthly', 'semiannual', 'annual');

-- DropForeignKey
ALTER TABLE "friends" DROP CONSTRAINT "friends_user_one_id_fkey";

-- DropForeignKey
ALTER TABLE "friends" DROP CONSTRAINT "friends_user_second_id_fkey";

-- DropForeignKey
ALTER TABLE "song_of_the_day" DROP CONSTRAINT "song_of_the_day_user_id_fkey";

-- DropIndex
DROP INDEX "song_of_the_day_trackId_key";

-- AlterTable
ALTER TABLE "song_of_the_day" ADD COLUMN     "mood" TEXT,
ADD COLUMN     "note" TEXT;

-- DropTable
DROP TABLE "friends";

-- DropTable
DROP TABLE "users";

-- CreateTable
CREATE TABLE "user" (
    "id" TEXT NOT NULL,
    "username" TEXT NOT NULL,
    "passwordHash" TEXT NOT NULL,
    "is_verified" BOOLEAN NOT NULL DEFAULT false,
    "email" TEXT NOT NULL,
    "bio" TEXT,
    "photo_url" TEXT,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "user_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "friend" (
    "id" TEXT NOT NULL,
    "status" "FriendStatus" NOT NULL,
    "user_one_id" TEXT NOT NULL,
    "user_two_id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "friend_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "user_statistic" (
    "id" TEXT NOT NULL,
    "user_id" TEXT NOT NULL,
    "period" "StatisticPeriod" NOT NULL,
    "totalTracks" INTEGER NOT NULL DEFAULT 0,
    "totalDuration" INTEGER NOT NULL DEFAULT 0,
    "uniqueArtists" INTEGER NOT NULL DEFAULT 0,
    "vibe" TEXT,
    "topArtistsIds" TEXT[],
    "topTracksIds" TEXT[],
    "topAlbumsIds" TEXT[],
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "user_statistic_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "user_username_key" ON "user"("username");

-- CreateIndex
CREATE UNIQUE INDEX "user_passwordHash_key" ON "user"("passwordHash");

-- CreateIndex
CREATE UNIQUE INDEX "user_email_key" ON "user"("email");

-- AddForeignKey
ALTER TABLE "friend" ADD CONSTRAINT "friend_user_one_id_fkey" FOREIGN KEY ("user_one_id") REFERENCES "user"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "friend" ADD CONSTRAINT "friend_user_two_id_fkey" FOREIGN KEY ("user_two_id") REFERENCES "user"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "song_of_the_day" ADD CONSTRAINT "song_of_the_day_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "user_statistic" ADD CONSTRAINT "user_statistic_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
