import dynamic from "next/dynamic";
import React from "react";

const Map = dynamic(() => import("@/app/components/Map/Map"), { ssr: false });

interface Sensor {
  sensor: string;
  coord: string;
}

interface Sector {
  coords: string;
  sensors: Sensor[];
}

const Page = async () => {
  const sectors = {
    "sector-server-a": {
      coords:
        "-34.91739651002616, -56.16210771871581;-34.91520592550176, -56.15995194170265;-34.915812960803784, -56.15930842617635;-34.91799473178162, -56.161485653707054",
      sensors: [
        {
          sensor: "sensor-1a",
          coord: "-34.91689987581042, -56.16092258868043",
        },
        {
          sensor: "sensor-2a",
          coord: "-34.91689987581040, -56.16092258868040",
        },
      ],
    },
    "sector-server-b": {
      coords:
        "-34.918423283163776,-56.161186695098884; -34.91810633129162,-56.161422729492195;  -34.91592285180286,-56.159191131591804; -34.91659198878654,-56.1585259437561",
    },
  };

  return (
    <div className="w-full h-full">
      <Map sectors={sectors} />
    </div>
  );
};

export default Page;
