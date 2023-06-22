import L from "leaflet";

const iconSensor = L.Icon.extend({
  options: {
    iconUrl: "https://leafletjs.com/examples/custom-icons/leaf-green.png",
    // iconRetinaUrl: require("/img/iconSensor.svg"),
    iconAnchor: null,
    popupAnchor: null,
    shadowUrl: null,
    shadowSize: null,
    shadowAnchor: null,
    iconSize: new L.Point(60, 75),
    className: "leaflet-div-icon",
  },
});

export { iconSensor };
