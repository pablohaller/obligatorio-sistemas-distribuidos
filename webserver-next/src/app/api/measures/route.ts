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