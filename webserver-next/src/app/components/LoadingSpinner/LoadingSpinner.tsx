import { IconLoader3 } from "@tabler/icons-react";
import React from "react";

const LoadingSpinner = () => (
  <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-10 flex flex-col items-center">
    <IconLoader3 className="h-20 w-20 text-sky-500 animate-spin" />
    <div className="text-3xl">Cargando</div>
  </div>
);

export default LoadingSpinner;
