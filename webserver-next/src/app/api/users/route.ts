import client from "@/app/prisma/client";
import { NextResponse } from "next/server";

export async function GET() {
  const users = await client.users.findMany();
  return NextResponse.json(users);
}
