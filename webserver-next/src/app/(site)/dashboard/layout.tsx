import { authOptions } from "@/app/api/auth/[...nextauth]/route";
import SideBar from "@/app/components/SideBar/SideBar";
import { getServerSession } from "next-auth";
import { headers } from "next/headers";
import React from "react";

interface Props {
  children: React.ReactNode;
}

const Layout = async ({ children }: Props) => {
  const session = await getServerSession(authOptions);
  const headersList = headers();
  const pathname = headersList.get("pathname") || "";

  return (
    <div className="flex flex-col md:flex-row h-full">
      <SideBar session={session} />
      <div>{children}</div>
    </div>
  );
};

export default Layout;
