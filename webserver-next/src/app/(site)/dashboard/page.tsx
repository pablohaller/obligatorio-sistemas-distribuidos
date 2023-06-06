import { authOptions } from "@/app/api/auth/[...nextauth]/route";
import Modal from "@/app/components/modal/modal";
import { getServerSession } from "next-auth";
import { useSession, signOut } from "next-auth/react";

const Page = async () => {
  return <>{/* {true && <Modal />} */}</>;
};

export default Page;
