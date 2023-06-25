import prisma from "@/app/prisma/client";
import bcrypt from "bcrypt";
import { NextRequest, NextResponse } from "next/server";

export async function GET() {
  const users = await prisma.users.findMany();
  return NextResponse.json(users);
}

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    const { name, email, password } = body;

    if (!name || !email || !password) {
      return NextResponse.json({ error: "Faltan campos" }, { status: 400 });
    }

    const exist = await prisma.users.findFirst({
      where: {
        email,
      },
    });

    if (exist) {
      return NextResponse.json(
        { error: "Cuenta de e-mail ya registrada" },
        {
          status: 400,
        }
      );
    }

    const hashedPassword = await bcrypt.hash(password, 10);

    const user = await prisma.users.create({
      data: {
        email,
        name,
        password: hashedPassword,
      },
    });

    return NextResponse.json(user, { status: 200 });
  } catch (error) {
    return NextResponse.json(
      { error: "Algo fue mal. Por favor, intente de nuevo" },
      { status: 500 }
    );
  }
}
