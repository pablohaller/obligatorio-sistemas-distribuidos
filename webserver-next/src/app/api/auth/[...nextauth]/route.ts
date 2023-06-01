import NextAuth from "next-auth/next";
import CredentialsProvider from "next-auth/providers/credentials";
import client from "@/app/prisma/client";
import bcrypt from "bcrypt";
import { PrismaAdapter } from "@next-auth/prisma-adapter";
import { SessionStrategy } from "next-auth";

const jwt: SessionStrategy = "jwt";

export const authOptions = {
  adapter: PrismaAdapter(client),
  providers: [
    CredentialsProvider({
      name: "credentials",
      credentials: {
        email: { label: "Email", type: "text", placeholder: "jsmith" },
        password: { label: "Password", type: "password" },
        name: {
          label: "Username",
          type: "text",
          placeholder: "John Smith",
        },
      },
      async authorize(credentials) {
        if (!credentials?.email || !credentials?.password) {
          throw new Error("Please enter an email and password");
        }

        const user = await client.users.findFirst({
          where: {
            email: credentials?.email,
          },
        });

        if (!user || !user?.password) {
          throw new Error("No user found");
        }

        const passwordMatch = await bcrypt.compare(
          credentials.password,
          user.password
        );

        if (!passwordMatch) {
          throw new Error("Incorrect password");
        }

        console.log(user);

        return user;
      },
    }),
  ],
  secret: process.env.NEXTAUTH_SECRET,
  session: {
    strategy: jwt,
  },
};

const handler = NextAuth(authOptions);

export { handler as GET, handler as POST };
