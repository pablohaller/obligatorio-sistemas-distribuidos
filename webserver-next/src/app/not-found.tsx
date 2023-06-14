import React from "react";
import LoadingSpinner from "./components/LoadingSpinner/LoadingSpinner";
import { IconCircleCheck, IconCircleX } from "@tabler/icons-react";
import Link from "next/link";

const Page = () => {
  return (
    <div className="h-screen grid place-items-center">
      <div className="   bg-white p-4 rounded-xl drop-shadow-md md:w-2/5 relative h-1/2">
        <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-10 flex flex-col items-center text-center">
          <IconCircleX className="h-20 w-20 text-red-500" />
          <div className="text-3xl">Página no encontrada</div>
          <Link className=" text-sky-500 hover:underline mt-2" href="/">
            Presione aquí para regresar al dashboard
          </Link>
        </div>
      </div>
    </div>
  );
};

export default Page;
