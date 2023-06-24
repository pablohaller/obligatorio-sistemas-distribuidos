import dynamic from "next/dynamic";
import { DateTime } from "luxon";
import { MeasureChartData } from "@/app/components/MeasureChart/MeasureChart";

const Map = dynamic(() => import("@/app/components/Map/Map"), { ssr: false });
const MeasureChart = dynamic(
  () => import("@/app/components/MeasureChart/MeasureChart"),
  { ssr: false }
);

interface Props {
  params: { id: string };
}

export const getData = async (payload: any) => {
  const request = await fetch(`${process.env.NEXT_NGINX_API_URI}/Alert`, {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });
  return request.json();
};

export const getMapData = async (payload: any) => {
  const request = await fetch(`${process.env.NEXT_NGINX_API_URI}/MapReport`, {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });
  return request.json();
};

const Page = async ({ params }: Props) => {
  const measure = await client.measures.findUnique({
    where: {
      id: params?.id,
    },
  });
  const parsedData = JSON.parse(measure?.data?.replace(/'/g, '"') || "");
  const measurements = await getData(parsedData);
  const mapData = await getMapData(parsedData);
  const sectors = { sector: parsedData?.sector, ...mapData };
  const mappedMeasurements = measurements.map(
    (measurement: MeasureChartData) => ({
      ...measurement,
      datetime: DateTime.fromISO(measurement?.datetime)
        .setZone("America/Montevideo")
        .toFormat("hh:mm:ss a"),
    })
  );

  return (
    <div className="w-full h-full">
      <div className="w-full h-2/3">
        <Map sectors={sectors} mapReport />
      </div>
      <div className="w-full h-1/3">
        <MeasureChart data={mappedMeasurements} />
      </div>
    </div>
  );
};

export default Page;
