// @ts-nocheck
"use client";
import {
  MapContainer,
  Marker,
  Polygon,
  Popup,
  TileLayer,
  useMap,
  useMapEvents,
} from "react-leaflet";
import "leaflet/dist/leaflet.css";
import MapData from "./MapData";

const Map = () => {
  const greenOptions = { color: "green" };
  const green = [
    [-34.91739651002616, -56.16210771871581],
    [-34.91520592550176, -56.15995194170265],
    [-34.915812960803784, -56.15930842617635],
    [-34.91799473178162, -56.161485653707054],
  ];

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
      <Polygon pathOptions={greenOptions} positions={green} />
      {/* <Marker position={[51.505, -0.09]}>
        <Popup>Test</Popup>
      </Marker> */}
    </MapContainer>
  );
};

export default Map;
