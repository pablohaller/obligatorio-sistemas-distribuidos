import LoadingSpinner from "@/app/components/LoadingSpinner/LoadingSpinner";
import MeasuresList from "@/app/components/MeasuresList/MeasuresList";
import { IconReport } from "@tabler/icons-react";
import React from "react";

const Page = async () => {
  const measures = await client.measures.findMany();

  return (
    <>
      <div className="font-rubik flex text-lg items-center border-b border-b-gray-200 p-2 pb-4 mb-2 sticky top-0 bg-white">
        <IconReport className="h-6 w-6 mr-2 text-sky-500" />
        Reportes
      </div>
      <MeasuresList measures={measures} />
    </>
  );
};

export default Page;
