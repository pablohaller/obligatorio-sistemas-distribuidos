"use client";
import React, { useEffect, useRef, useState } from "react";
import { DateTime } from "luxon";
import {
  CartesianGrid,
  Legend,
  Line,
  LineChart,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import SectorListItem from "./SectorListItem";

interface Props {
  sectors: any[];
}

const mockedData = {
  "sensor-1a": [
    {
      datetime: "2023-06-16T16:39:59Z",
      pressure: 112,
    },
    {
      datetime: "2023-06-16T16:40:09Z",
      pressure: 112,
    },
    {
      datetime: "2023-06-16T16:40:19Z",
      pressure: 60,
    },
    {
      datetime: "2023-06-16T16:40:29Z",
      pressure: 60,
    },
    {
      datetime: "2023-06-16T16:40:39Z",
      pressure: 112,
    },
    {
      datetime: "2023-06-16T16:40:49Z",
      pressure: 120,
    },
  ],
  "sensor-2a": [
    {
      datetime: "2023-06-16T16:39:59Z",
      pressure: 112,
    },
    {
      datetime: "2023-06-16T16:40:09Z",
      pressure: 112,
    },
    {
      datetime: "2023-06-16T16:40:19Z",
      pressure: 112,
    },
    {
      datetime: "2023-06-16T16:40:29Z",
      pressure: 112,
    },
    {
      datetime: "2023-06-16T16:40:39Z",
      pressure: 112,
    },
    {
      datetime: "2023-06-16T16:40:49Z",
      pressure: 112,
    },
  ],
  "sensor-3a": [
    {
      datetime: "2023-06-16T16:39:59Z",
      pressure: 112,
    },
    {
      datetime: "2023-06-16T16:40:09Z",
      pressure: 112,
    },
    {
      datetime: "2023-06-16T16:40:19Z",
      pressure: 112,
    },
    {
      datetime: "2023-06-16T16:40:29Z",
      pressure: 112,
    },
    {
      datetime: "2023-06-16T16:40:39Z",
      pressure: 112,
    },
    {
      datetime: "2023-06-16T16:40:49Z",
      pressure: 112,
    },
  ],
  "sensor-4a": [
    {
      datetime: "2023-06-16T16:39:59Z",
      pressure: 112,
    },
    {
      datetime: "2023-06-16T16:40:09Z",
      pressure: 112,
    },
    {
      datetime: "2023-06-16T16:40:19Z",
      pressure: 112,
    },
    {
      datetime: "2023-06-16T16:40:29Z",
      pressure: 112,
    },
    {
      datetime: "2023-06-16T16:40:39Z",
      pressure: 112,
    },
    {
      datetime: "2023-06-16T16:40:49Z",
      pressure: 112,
    },
  ],
};

const SectorList = ({ sectors }: Props) => {
  const [chartData, setChartData] = useState([]);
  const [selectedSector, setSelectedSector] = useState(
    sectors?.[0]?.sector || null
  );
  const chartsDiv = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const getData = async () => {
      // const response = await fetch(
      //   `${process.env.NEXT_PUBLIC_NGINX_API_URI}/LastSectorMeasurements/${selectedSector}/1` ||
      //     ""
      // );
      // const data = await response.json();
      // console.log("data", data);
      const dataArray = Object.keys(mockedData)?.reduce((array: any, key) => {
        // return [...array, { sensor: key }];
        const data = mockedData[key as keyof typeof mockedData];
        return [
          ...array,
          // { sensor: key, data: mockedData[key as keyof typeof mockedData] },
          {
            sensor: key,
            data: data?.map((obj) => ({
              ...obj,
              pressure: obj?.pressure * (Math.floor(Math.random() * 10) + 1),
              datetime: DateTime.fromISO(obj?.datetime)
                .setZone("America/Montevideo")
                .toFormat("dd/MM/yyyy - hh:mm:ss a"),
            })),
          },
        ];
      }, []);
      setChartData(dataArray);
    };
    getData();
  }, [selectedSector]);

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
        <span>Mediciones del sector&nbsp;</span>
        <span className="text-sky-500">{selectedSector}</span>:
      </div>
      <div className="p-2 h-3/4 overflow-auto" ref={chartsDiv}>
        <div className="w-full h-full grid md:grid-cols-2 gap-2">
          {selectedSector &&
            !!chartData?.length &&
            chartData?.map(({ sensor, data }) => (
              <div key={`sensor-${sensor}`}>
                <div className="py-2 font-semibold text-sky-500">{sensor}</div>
                <div className="w-full">
                  <ResponsiveContainer width="100%" height={200}>
                    <LineChart
                      data={data}
                      margin={{
                        top: 5,
                        right: 30,
                        left: 20,
                        bottom: 5,
                      }}
                    >
                      <Tooltip />
                      <CartesianGrid strokeDasharray="3 3" />
                      <XAxis dataKey="datetime" />
                      <YAxis />
                      <Legend />
                      <Line
                        name="Presión"
                        type="monotone"
                        dataKey="pressure"
                        stroke="#0BA5E9"
                        activeDot={{ r: 8 }}
                      />
                    </LineChart>
                  </ResponsiveContainer>
                </div>
              </div>
            ))}
        </div>
      </div>
    </div>
  );
};

export default SectorList;
