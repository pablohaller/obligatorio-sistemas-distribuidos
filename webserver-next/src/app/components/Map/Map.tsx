// @ts-nocheck
"use client";
import {
  MapContainer,
  Marker,
  Polygon,
  Popup,
  TileLayer,
  Tooltip,
} from "react-leaflet";
import "leaflet/dist/leaflet.css";
import MapData from "./MapData";
import { useMemo } from "react";
interface Props {
  sectors: unknown;
}

const Map = ({ sectors }: Props) => {
  const mapData = useMemo(() => {
    let mapSectors = [];
    let mapSensors = [];
    Object.keys(sectors).forEach((key) => {
      const sector = sectors[key];
      sector?.sensors?.forEach((sensor) => {
        mapSensors.push({
          sensor: key,
          coords: sensor?.coord?.split(",")?.map((unit) => parseFloat(unit)),
        });
      });
      mapSectors.push({
        sector: key,
        sensors: sector?.sensors,
        positions: sector?.coords
          .split(";")
          ?.map((coord) => coord.split(",")?.map((unit) => parseFloat(unit))),
        pathOptions: { color: "green" },
      });
    });

    return {
      mapSectors,
      mapSensors,
    };
  }, [sectors]);

  const { mapSectors, mapSensors } = mapData;

  return (
    <MapContainer
      center={[-34.87156680929775, -56.2262692906801]}
      zoom={12}
      scrollWheelZoom={true}
    >
      <TileLayer
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />
      <MapData />
      {mapSectors?.map(({ sector, positions, sensors, pathOptions }) => (
        <Polygon
          key={`sector-polygon-${sector}`}
          positions={positions}
          pathOptions={pathOptions}
        >
          <Tooltip>{sector}</Tooltip>
          <Popup>
            {sensors?.map(({ sensor }) => (
              <div key={`popup-${sector}-sensor-${sensor}`}>{sensor}</div>
            ))}
          </Popup>
        </Polygon>
      ))}
      {mapSensors?.map(({ sensor, coords }) => (
        <Marker
          key={`sensor-marker-${sensor}`}
          position={coords}
          icon={L.icon({
            iconUrl: "/img/iconSensor.png",
            iconSize: [44, 44],
          })}
        >
          <Popup>{sensor}</Popup>
          <Tooltip>{sensor}</Tooltip>
        </Marker>
      ))}
    </MapContainer>
  );
};

export default Map;
