import { DASHBOARD_REPORTS_PATH } from "@/app/constants/routes";
import prisma from "@/app/prisma/client";
import { Users } from "@prisma/client";
import { NextResponse } from "next/server";
import nodemailer from "nodemailer";

export async function GET() {
  return NextResponse.json({ data: "Pepe" });
}

export async function POST(request: Request) {
  const { data, filtration } = await request.json();
  const savedMeasure = await prisma.measures.create({
    data: {
      data,
      filtration,
    },
  });

  if (filtration) {
    try {
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
        subject: "Filtración encontrada",
        html: `<a href="${process.env.NEXT_API_URL}${DASHBOARD_REPORTS_PATH}/${savedMeasure?.id}">Click para ver</a>`,
      };

      transporter.sendMail(mailOptions, (error, info) => {
        if (error) {
          console.error(error);
        } else {
          console.log("Email sent: " + info.response);
        }
      });
    } catch (error) {
      console.error("Email not sent", error);
    }
  }

  return NextResponse.json(savedMeasure);
}

export async function DELETE(request: Request) {
  const { searchParams } = new URL(request.url);
  const id = searchParams.get("id");
  let response: any = {};
  if (id) {
    const measure = await prisma.measures.findUnique({
      where: {
        id,
      },
    });
    if (measure) {
      if (measure?.isActive) {
        const updatedMeasure = await prisma.measures.update({
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
