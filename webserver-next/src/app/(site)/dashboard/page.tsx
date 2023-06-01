"use client";
import { useSession, signOut } from "next-auth/react";

const Page = () => {
  const session = useSession();

  return (
    <div>
      <h1>Dashboard</h1>
      <p>{JSON.stringify(session)}</p>
      <button onClick={() => signOut()}>Sign Out</button>
    </div>
  );
};

export default Page;
