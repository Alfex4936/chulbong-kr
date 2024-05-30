import MyLocateOverlay from "@/components/map/MyLocateOverlay";
import useMapStatusStore from "@/store/useMapStatusStore";
import useMapStore from "@/store/useMapStore";
import { createRoot } from "react-dom/client";

const useMyGps = () => {
  const { setPosition } = useMapStatusStore();
  const { map, myLocateOverlay, setMyLocateOverlay } = useMapStore();

  const centerMapOnCurrentPosition = () => {
    if (map && navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const moveLatLon = new window.kakao.maps.LatLng(
            position.coords.latitude,
            position.coords.longitude
          );

          if (myLocateOverlay) {
            myLocateOverlay.setMap(null);
          }

          const overlayDiv = document.createElement("div");
          const root = createRoot(overlayDiv);
          root.render(<MyLocateOverlay />);

          const customOverlay = new window.kakao.maps.CustomOverlay({
            position: moveLatLon,
            content: overlayDiv,
            zIndex: 10,
          });

          customOverlay.setMap(map);
          setMyLocateOverlay(customOverlay);

          setPosition(position.coords.latitude, position.coords.longitude);
          map.setCenter(moveLatLon);
        },
        () => {
          alert("잠시 후 다시 시도해주세요.");
        }
      );
    } else {
      alert("잠시 후 다시 시도해주세요.");
    }
  };
  const centerMapOnCurrentPositionAsync = (callback: any) => {
    if (map && navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        async (position) => {
          const moveLatLon = new window.kakao.maps.LatLng(
            position.coords.latitude,
            position.coords.longitude
          );

          if (myLocateOverlay) {
            myLocateOverlay.setMap(null);
          }

          const overlayDiv = document.createElement("div");
          const root = createRoot(overlayDiv);
          root.render(<MyLocateOverlay />);

          const customOverlay = new window.kakao.maps.CustomOverlay({
            position: moveLatLon,
            content: overlayDiv,
            zIndex: 10,
          });

          customOverlay.setMap(map);
          setMyLocateOverlay(customOverlay);

          const move = async () => {
            setPosition(position.coords.latitude, position.coords.longitude);
            map.setCenter(moveLatLon);
          };

          await move();

          callback();
        },
        () => {
          alert("잠시 후 다시 시도해주세요.");
        }
      );
    } else {
      alert("잠시 후 다시 시도해주세요.");
    }
  };

  return { centerMapOnCurrentPosition, centerMapOnCurrentPositionAsync };
};

export default useMyGps;
