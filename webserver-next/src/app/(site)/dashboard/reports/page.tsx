import LoadingSpinner from "@/app/components/LoadingSpinner/LoadingSpinner";
import MeasuresList from "@/app/components/MeasuresList/MeasuresList";
import { IconReport } from "@tabler/icons-react";
import React from "react";

const Page = async () => {
  const measures = await client.measures.findMany();

  return (
    <div>
      {/* To Component Page Header */}
      <div className="font-rubik flex text-4xl items-center border-b border-b-gray-200 p-2 pb-4 mb-2">
        <IconReport className="h-10 w-10 mr-2 text-sky-500" />
        Reportes
      </div>
      <MeasuresList measures={measures} />
    </div>
  );
};

export default Page;
