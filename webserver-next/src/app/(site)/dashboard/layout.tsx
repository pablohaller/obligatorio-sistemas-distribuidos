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
      <div className="w-full h-screen p-2">
        <div className="bg-white/80 h-full p-2 rounded-xl shadow-xl border border-gray-200 relative break-all overflow-y-scroll">
          {children}
        </div>
      </div>
    </div>
  );
};

export default Layout;
