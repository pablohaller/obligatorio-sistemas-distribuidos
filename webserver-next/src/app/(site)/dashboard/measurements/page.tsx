import { IconEngine } from "@tabler/icons-react";
import dynamic from "next/dynamic";
import React, { useState } from "react";

const SectorList = dynamic(
  () => import("@/app/components/SectorList/SectorList"),
  { ssr: false }
);

async function getData() {
  const res = await fetch(`${process.env.NEXT_NGINX_API_URI}/Sectors` || "");
  if (!res.ok) {
    throw new Error("Failed to fetch data");
  }

  return res.json();
}

const Page = async () => {
  const sectors = await getData();

  return (
    <div className="h-[75%]">
      <div className="font-rubik flex text-lg items-center border-b border-b-gray-200 p-2 pb-4 mb-2 sticky top-0 bg-white">
        <IconEngine className="h-6 w-6 mr-2 text-sky-500" />
        Mediciones
      </div>
      <div className="h-full">
        <SectorList sectors={sectors} />
      </div>
    </div>
  );
};

export default Page;
