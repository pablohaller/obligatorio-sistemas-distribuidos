import { PrismaClient } from "@prisma/client";

let prisma: PrismaClient;

if (process.env.NODE_ENV === "production") {
  prisma = new PrismaClient();
} else {
  if (!global?.client) {
    global.client = new PrismaClient();
  }

  prisma = global?.client;
}

export default prisma;
