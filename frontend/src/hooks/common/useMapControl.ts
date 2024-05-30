import { MAP_LAT_DIF } from "@/constants";
import useMapStatusStore from "@/store/useMapStatusStore";
import useMapStore from "@/store/useMapStore";
import useMiniMapStatusStore from "@/store/useMiniMapStatusStore";
import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";

const useMapControl = () => {
  const { map, markers } = useMapStore();
  const { setPosition } = useMapStatusStore();
  const { open } = useMobileMapOpenStore();
  const { map: miniMap } = useMiniMapStatusStore();

  const filterMarker = async (markerId: number) => {
    const imageSize = new window.kakao.maps.Size(39, 39);
    const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

    const selectedMarkerImg = new window.kakao.maps.MarkerImage(
      "/selectedMarker.svg",
      imageSize,
      imageOption
    );

    const activeMarkerImg = new window.kakao.maps.MarkerImage(
      "/activeMarker.svg",
      imageSize,
      imageOption
    );

    markers?.forEach((marker) => {
      if (Number(marker.getTitle()) === markerId) {
        marker.setImage(selectedMarkerImg);
      } else {
        marker.setImage(activeMarkerImg);
      }
    });
  };

  const moveLocation = async (lat: number, lng: number, mini?: boolean) => {
    if (mini) {
      const moveLatLon = new window.kakao.maps.LatLng(
        (lat as number) + MAP_LAT_DIF,
        lng
      );

      miniMap?.panTo(moveLatLon);
    } else {
      const moveLatLon = new window.kakao.maps.LatLng(
        (lat as number) + MAP_LAT_DIF,
        lng
      );
      setPosition((lat as number) + MAP_LAT_DIF, lng as number);
      map?.panTo(moveLatLon);
      open();
    }
  };

  return { filterMarker, moveLocation };
};

export default useMapControl;
