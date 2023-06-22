"use client";
import { useEffect } from "react";
import { useMap, useMapEvent, useMapEvents } from "react-leaflet";

const MapData = () => {
  const map = useMapEvents({
    drag: (location) => {
      console.log(map.getCenter());
    },
    click: (location) => {
      console.log(location.latlng);
    },
    load: () => {
      console.log("loaded");
    },
    // click: () => {
    //   map.locate();
    // },
    // locationfound: (location) => {
    //   console.log("location found:", location);
    // },
  });
  return null;
};

export default MapData;
