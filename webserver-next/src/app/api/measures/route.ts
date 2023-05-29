import client from "@/app/prisma/client";
import { NextResponse } from "next/server";

export async function GET() {
  return NextResponse.json({ data: "Pepe" });
}

export async function POST(request: Request) {
  const { data } = await request.json();
  const savedMeasure = await client.measures.create({
    data: {
      data,
    },
  });
  return NextResponse.json(savedMeasure);
}
