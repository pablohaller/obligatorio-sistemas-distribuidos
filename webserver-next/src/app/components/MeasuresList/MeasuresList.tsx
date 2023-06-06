"use client";
import { Measures } from "@prisma/client";
import { IconEngine, IconEngineOff } from "@tabler/icons-react";
import { DateTime } from "luxon";
import React, { useEffect, useState } from "react";
import { toast } from "react-toastify";

interface Props {
  measures: Measures[];
}

const MeasuresList = ({ measures }: Props) => {
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<boolean>(false);

  return (
    <div>
      {measures?.map(({ id, data, filtration }) => {
        const parsedData = JSON.parse(data.replace(/'/g, '"'));
        console.log(parsedData);
        return (
          <div
            key={id}
            className="border rounded-xl p-2 flex border-l-8 border-l-red-500 flex-col md:flex-row"
          >
            <div className="flex-shrink-0 m-2">
              {filtration ? (
                <IconEngineOff className="h-10 w-10 text-red-400" />
              ) : (
                <IconEngine className="h-10 w-10" />
              )}
            </div>
            <div>
              <div className="border-b border-b-gray-200 pb-2">
                <span className="font-rubik  mr-2">Reporte:</span>
                <span>{id}</span>
              </div>
              <div>
                {parsedData?.map((measure: any, index: number) => (
                  <div
                    key={`${id}/${index}`}
                    className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-2"
                  >
                    <div>
                      <div className="font-rubik break-normal">
                        Fecha y hora
                      </div>
                      <div className="break-normal text-sm">
                        {DateTime.fromISO(measure?.Datetime)
                          .setZone("America/Montevideo")
                          .toFormat("dd/MM/yyyy - hh:mm a")}
                      </div>
                    </div>
                    <div>
                      <div className="font-rubik">Presi&oacute;n</div>
                      <div className="break-normal text-sm">
                        {measure?.Presion}
                      </div>
                    </div>
                    <div>
                      <div className="font-rubik">Sector</div>
                      <div className="break-normal text-sm">
                        {measure?.Sector}
                      </div>
                    </div>
                    <div>
                      <div className="font-rubik">Sensor</div>
                      <div className="break-normal text-sm">
                        {measure?.Sensor}
                      </div>
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
