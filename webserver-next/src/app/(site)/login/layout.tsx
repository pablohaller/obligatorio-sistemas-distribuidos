import { authOptions } from "@/app/api/auth/[...nextauth]/route";
import { DASHBOARD_PATH } from "@/app/constants/routes";
import { getServerSession } from "next-auth";
import { redirect } from "next/navigation";

interface Props {
  children: React.ReactNode;
}

const Layout = async ({ children }: Props) => {
  const session = await getServerSession(authOptions);
  if (session) {
    redirect(DASHBOARD_PATH);
  }
  return <div>{children}</div>;
};

export default Layout;
