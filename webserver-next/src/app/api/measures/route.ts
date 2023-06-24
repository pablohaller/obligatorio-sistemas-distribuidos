import client from "@/app/prisma/client";
import { Users } from "@prisma/client";
import { NextResponse } from "next/server";
import nodemailer from "nodemailer";

export async function GET() {
  return NextResponse.json({ data: "Pepe" });
}

export async function POST(request: Request) {
  const { data, filtration } = await request.json();
  const savedMeasure = await client.measures.create({
    data: {
      data,
      filtration,
    },
  });

  if (filtration) {
    const req = await fetch(`${process.env.NEXT_API_URL}/api/users`);
    const users = await req.json();
    const mailingList = users.map(({ email }: Users) => email).join(",");

    const transporter = nodemailer.createTransport({
      host: "smtp-mail.outlook.com", // hostname
      secure: false, // TLS requires secureConnection to be false
      port: 587, // port for secure SMTP
      tls: {
        ciphers: "SSLv3",
      },
      auth: {
        user: process.env.GMAIL_ACCOUNT,
        pass: process.env.GMAIL_PASSWORD,
      },
    });

    const mailOptions = {
      from: process.env.GMAIL_ACCOUNT,
      to: mailingList,
      subject: "Lorem Ipsum",
      text: "Lorem Ipsum",
    };

    transporter.sendMail(mailOptions, (error, info) => {
      if (error) {
        console.error(error);
      } else {
        console.log("Email sent: " + info.response);
      }
    });
  }

  return NextResponse.json(savedMeasure);
}

export async function DELETE(request: Request) {
  const { searchParams } = new URL(request.url);
  const id = searchParams.get("id");
  let response: any = {};
  if (id) {
    const measure = await client.measures.findUnique({
      where: {
        id,
      },
    });
    if (measure) {
      if (measure?.isActive) {
        const updatedMeasure = await client.measures.update({
          where: {
            id,
          },
          data: {
            isActive: false,
          },
        });
        response = updatedMeasure;
      } else {
        response = { error: "Measure already inactive" };
      }
    } else {
      response = { error: "Measure not found" };
    }
  } else {
    response = { error: "No Id Provided" };
  }

  return NextResponse.json(response, {
    ...(response?.error ? { status: 400 } : {}),
  });
}
