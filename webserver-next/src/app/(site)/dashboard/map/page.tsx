import dynamic from "next/dynamic";

const Map = dynamic(() => import("@/app/components/Map/Map"), { ssr: false });

interface Sensor {
  sensor: string;
  coord: string;
}

interface Sector {
  coords: string;
  sensors: Sensor[];
}

async function getData() {
  const res = await fetch(`${process.env.NEXT_NGINX_API_URI}/Map` || "");
  if (!res.ok) {
    throw new Error("Failed to fetch data");
  }

  return res.json();
}

const Page = async () => {
  const sectors = await getData();

  return (
    <div className="w-full h-full">
      <Map sectors={sectors} />
    </div>
  );
};

export default Page;
