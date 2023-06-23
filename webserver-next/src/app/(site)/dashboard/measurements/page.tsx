import { IconEngine } from "@tabler/icons-react";
import dynamic from "next/dynamic";
import React, { useState } from "react";

const SectorList = dynamic(
  () => import("@/app/components/SectorList/SectorList"),
  { ssr: false }
);

const Page = async () => {
  // const sectors = await fetch('GET SECTORS');
  const sectors = [
    "sector-server-a-1",
    "sector-server-a-2",
    "sector-server-a-3",
    "sector-server-a-4",
    "sector-server-a-5",
    "sector-server-a-6",
    "sector-server-a-7",
    "sector-server-a-8",
    "sector-server-a-9",
    "sector-server-a-10",
  ];

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
