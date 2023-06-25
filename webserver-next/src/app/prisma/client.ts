import { PrismaClient } from "@prisma/client";

let prisma: PrismaClient;

if (process.env.NODE_ENV === "production") {
  prisma = new PrismaClient();
  // TRNCATE REPORTS
  const deleteAll = async () => {
    const deleteReports = await prisma.measures.deleteMany({ where: {} });
    console.log("Delete...", deleteReports);
  };
  deleteAll();
} else {
  if (!global?.client) {
    global.client = new PrismaClient();
  }

  prisma = global?.client;
}

export default prisma;
