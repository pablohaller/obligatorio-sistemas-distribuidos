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

      const html = `<link rel="preconnect" href="https://fonts.googleapis.com" />
      <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
      <link
        href="https://fonts.googleapis.com/css2?family=Rubik:wght@300;400&display=swap"
        rel="stylesheet"
      />
      <div
        style="
          border: 1px solid #0ba5e9;
          font-family: Rubik;
          width: 400px;
          padding: 1rem;
          text-align: center;
        "
      >
        <div style="margin-bottom: 1rem;">
            <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-engine-off" width="44" height="44" viewBox="0 0 24 24" stroke-width="1.5" stroke="#ff2825" fill="none" stroke-linecap="round" stroke-linejoin="round">
              <path stroke="none" d="M0 0h24v24H0z" fill="none"/>
              <path d="M3 10v6" />
              <path d="M12 5v3" />
              <path d="M10 5h4" />
              <path d="M5 13h-2" />
              <path d="M16 16h-1v2a1 1 0 0 1 -1 1h-3.465a1 1 0 0 1 -.832 -.445l-1.703 -2.555h-2v-6h2l.99 -.99m3.01 -1.01h1.382a1 1 0 0 1 .894 .553l1.448 2.894a1 1 0 0 0 .894 .553h1.382v-2h2a1 1 0 0 1 1 1v6" />
              <path d="M3 3l18 18" />
            </svg>
            <path stroke="none" d="M0 0h24v24H0z" fill="none" />
            <path d="M3 10v6" />
            <path d="M12 5v3" />
            <path d="M10 5h4" />
            <path d="M5 13h-2" />
            <path
              d="M16 16h-1v2a1 1 0 0 1 -1 1h-3.465a1 1 0 0 1 -.832 -.445l-1.703 -2.555h-2v-6h2l.99 -.99m3.01 -1.01h1.382a1 1 0 0 1 .894 .553l1.448 2.894a1 1 0 0 0 .894 .553h1.382v-2h2a1 1 0 0 1 1 1v6"
            />
            <path d="M3 3l18 18" />
          </svg>
        </div>
        <div style="font-weight: 400; margin-bottom: 1rem;">
          Se ha detectado una falla en un sector.
        </div>
        <div style="font-weight: 400; margin-bottom: 1rem;">
          Por favor, ingrese al siguiente <a href="${process.env.NEXT_API_URL}${DASHBOARD_REPORTS_PATH}/${savedMeasure?.id}">enlace</a> para obtener más
          información de la misma.
        </div>
        <div style="font-weight: 400; margin-bottom: 1rem;">Muchas gracias.</div>
      </div>
      `;

      const mailOptions = {
        from: process.env.GMAIL_ACCOUNT,
        to: mailingList,
        subject: "Filtración encontrada",
        html,
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

export async function PUT(request: Request) {
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

export async function DELETE() {
  // TRNCATE REPORTS
  const deleteAll = async () => {
    const deleteReports = await prisma.measures.deleteMany({ where: {isActive : true} });
    console.log("Delete...", deleteReports);
  };
  deleteAll();
  return NextResponse.json({ data: "ok" });
}