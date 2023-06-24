"use client";
import React, { useEffect, useRef, useState } from "react";
import { DateTime } from "luxon";
import SectorListItem from "./SectorListItem";
import LoadingSpinner from "../LoadingSpinner/LoadingSpinner";
import { IconRefresh } from "@tabler/icons-react";
import MeasureChart from "../MeasureChart/MeasureChart";

interface Props {
  sectors: any[];
}

const SectorList = ({ sectors }: Props) => {
  const [loading, setLoading] = useState(true);
  const [timespan, setTimespan] = useState(1);
  const [chartData, setChartData] = useState([]);
  const [selectedSector, setSelectedSector] = useState(
    sectors?.[0]?.sector || null
  );
  const chartsDiv = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const getData = async () => {
      console.log(
        "fetch",
        `${process.env.NEXT_PUBLIC_NGINX_API_URI}/LastSectorMeasurements/${selectedSector}/${timespan}`
      );
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_NGINX_API_URI}/LastSectorMeasurements/${selectedSector}/${timespan}` ||
          ""
      );
      const measurments = await response.json();
      setChartData(
        Object.keys(measurments)?.reduce(
          (array: any, key) => [
            ...array,
            {
              sensor: key,
              data: measurments[key as keyof typeof measurments]?.map(
                (obj: any) => ({
                  ...obj,
                  datetime: DateTime.fromISO(obj?.datetime)
                    .setZone("America/Montevideo")
                    .plus({ hours: -3 })
                    .toFormat("hh:mm:ss a"),
                })
              ),
            },
          ],
          []
        )
      );
      setLoading(false);
    };
    getData();
  }, [selectedSector, timespan]);

  if (loading) {
    return <LoadingSpinner />;
  }

  return (
    <div className="h-full">
      <div className="pl-2 font-semibold py-2">Sectores:</div>
      <div className="p-2 grid grid-cols-2 gap-2 h-1/4 overflow-auto overflow-y-scroll">
        {sectors?.map(({ sector }) => (
          <SectorListItem
            key={`sector-list-${sector}`}
            sector={sector}
            onClick={() => {
              setSelectedSector(sector);
              chartsDiv?.current?.scroll({
                top: 0,
                behavior: "smooth",
              });
            }}
            isSelected={selectedSector === sector}
          />
        ))}
      </div>
      <div className="pl-2 font-semibold py-2">
        <span title="Recargar">
          <IconRefresh className=" text-sky-500 hover:underline inline mr-2" />
        </span>
        <span>Mediciones de&nbsp;</span>
        <span className="text-sky-500 ">{selectedSector}&nbsp;</span>
        <span>hace &nbsp;</span>
        <select
          className="border border-sky-500 focus-visible:outline-sky-500 px-5 rounded-lg"
          onChange={(e) => setTimespan(Number(e.target.value))}
        >
          <option value="1">1 min</option>
          <option value="5">5 min</option>
          <option value="10">10 min</option>
        </select>
        <span>&nbsp;atrás:</span>
      </div>
      <div className="p-2 h-3/4 overflow-auto" ref={chartsDiv}>
        <div className="w-full h-full grid md:grid-cols-2 gap-2">
          {selectedSector &&
            !!chartData?.length &&
            chartData?.map(({ sensor, data }) => (
              <div key={`sensor-${sensor}`}>
                <div className="py-2 font-semibold text-sky-500">{sensor}</div>
                <div className="w-full">
                  <MeasureChart data={data} />
                </div>
              </div>
            ))}
        </div>
      </div>
    </div>
  );
};

export default SectorList;
