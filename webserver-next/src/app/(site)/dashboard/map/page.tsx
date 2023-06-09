import dynamic from "next/dynamic";
import React from "react";

const Map = dynamic(() => import("@/app/components/Map/Map"), { ssr: false });

const Page = () => {
  return (
    <div className="w-full h-full">
      <Map />
    </div>
  );
};

export default Page;
