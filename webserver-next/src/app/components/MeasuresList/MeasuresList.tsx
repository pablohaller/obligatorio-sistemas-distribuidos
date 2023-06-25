"use client";
import { DASHBOARD_REPORTS_PATH } from "@/app/constants/routes";
import { Measures } from "@prisma/client";
import {
  IconArrowRight,
  IconEngine,
  IconEngineOff,
  IconEye,
} from "@tabler/icons-react";
import { DateTime } from "luxon";
import Link from "next/link";
import { useState } from "react";

interface Props {
  measures: Measures[];
}

const MeasuresList = ({ measures }: Props) => {
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<boolean>(false);

  if (measures?.length === 0) {
    return (
      <div className="grid place-items-center w-full h-full">
        No hay reportes para mostrar
      </div>
    );
  }

  return (
    <div className="p-2">
      {measures?.map(({ id, data, filtration }) => {
        const parsedData = JSON.parse(data.replace(/'/g, '"'));
        return (
          <div
            key={id}
            className="border rounded-xl p-2 flex border-l-8 border-l-red-500 flex-col md:flex-row bg-white mb-2"
          >
            <div className="flex-shrink-0 m-2 md:mr-6 flex">
              {filtration ? (
                <IconEngineOff className="flex-shrink-0 h-10 w-10 text-red-400 mr-6 md:mr-0" />
              ) : (
                <IconEngine className="flex-shrink-0 h-10 w-10 mr-6 md:mr-0" />
              )}
              <div className="border-b border-b-gray-200 pb-2 md:hidden w-full">
                <div className="font-rubik  mr-2">Reporte:</div>
                <div>{id}</div>
              </div>
            </div>
            <div className="w-full">
              <div className="border-b border-b-gray-200 pb-2 hidden md:flex">
                <span className="font-rubik  mr-2">Reporte:</span>
                <span>{id}</span>
              </div>
              <div>
                {[parsedData]?.map((measure: any, index: number) => (
                  <div
                    key={`${id}/${index}`}
                    className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-2"
                  >
                    <div>
                      <div className="font-rubik break-normal">
                        Fecha y hora
                      </div>
                      <div className="break-normal text-sm">
                        {DateTime.fromISO(measure?.datetime)
                          .setZone("America/Montevideo")
                          .plus({ hours: -3 })
                          .toFormat("dd/MM/yyyy - hh:mm a")}
                      </div>
                    </div>
                    <div>
                      <div className="font-rubik">Presi&oacute;n</div>
                      <div className="break-normal text-sm">
                        {measure?.pressure}
                      </div>
                    </div>
                    <div>
                      <div className="font-rubik">Sector</div>
                      <div className="break-normal text-sm">
                        {measure?.sector}
                      </div>
                    </div>
                    <div>
                      <div className="font-rubik">Sensor</div>
                      <div className="break-normal text-sm">
                        {measure?.sensor}
                      </div>
                    </div>
                    <div className="flex justify-end w-full col-span-2 md:col-span-4">
                      <Link
                        className="text-sm flex items-center   hover:text-sky-500 cursor-pointer"
                        href={`${DASHBOARD_REPORTS_PATH}/${id}`}
                      >
                        <IconEye className="mr-2" />
                        <span className="hover:underline">Ver detalle</span>
                      </Link>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default MeasuresList;
