/*
  Warnings:

  - You are about to drop the column `published` on the `Measures` table. All the data in the column will be lost.

*/
-- RedefineTables
PRAGMA foreign_keys=OFF;
CREATE TABLE "new_Measures" (
    "id" TEXT NOT NULL PRIMARY KEY,
    "data" TEXT NOT NULL,
    "filtration" BOOLEAN NOT NULL DEFAULT false
);
INSERT INTO "new_Measures" ("data", "id") SELECT "data", "id" FROM "Measures";
DROP TABLE "Measures";
ALTER TABLE "new_Measures" RENAME TO "Measures";
PRAGMA foreign_key_check;
PRAGMA foreign_keys=ON;
