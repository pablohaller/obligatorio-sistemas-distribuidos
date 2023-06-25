import {
  DASHBOARD_MAP_PATH,
  DASHBOARD_MEASUREMENTS_PATH,
  DASHBOARD_REPORTS_PATH,
} from "@/app/constants/routes";
import { IconEngine, IconMap, IconReport } from "@tabler/icons-react";
import Link from "next/link";

const Page = async () => {
  return (
    <div className="w-full h-full grid grid-cols-1 md:grid-cols-2 gap-2 p-2 overflow-auto">
      <Link
        href={DASHBOARD_MAP_PATH}
        className="border rounded-xl p-2 flex-col bg-white border-l-8 border-l-sky-500 hover:bg-sky-100 grid place-items-center"
      >
        <div className="flex flex-col justify-center items-center">
          <IconMap className="h-14 w-14 text-sky-500" />
          <span className="text-3xl">Mapa</span>
          <span className="text-xs w-3/4 mt-4">
            Permite ver los sectores y sensores en un mapa de la ciudad
          </span>
        </div>
      </Link>
      <Link
        href={DASHBOARD_REPORTS_PATH}
        className="border rounded-xl p-2 flex-col bg-white border-l-8 border-l-sky-500 hover:bg-sky-100 grid place-items-center"
      >
        <div className="flex flex-col justify-center items-center">
          <IconReport className="h-14 w-14 text-sky-500" />
          <span className="text-3xl">Reportes</span>
          <span className="text-xs w-3/4 mt-4">
            Permite ver una lista con los últimos reportes de alerta de
            filtración si estos no han sido solucionados
          </span>
        </div>
      </Link>
      <Link
        href={DASHBOARD_MEASUREMENTS_PATH}
        className="border rounded-xl p-2 flex-col bg-white border-l-8 border-l-sky-500 hover:bg-sky-100 grid place-items-center"
      >
        <div className="flex flex-col justify-center items-center">
          <IconEngine className="h-14 w-14 text-sky-500" />
          <span className="text-3xl">Mediciones</span>
          <span className="text-xs w-3/4 mt-4">
            Permite ver una lista de sectores y las mediciones para cada sensor
            perteneciente a éste
          </span>
        </div>
      </Link>
    </div>
  );
};

export default Page;
