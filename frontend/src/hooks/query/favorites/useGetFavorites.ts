import { useQuery } from "@tanstack/react-query";
import getFavorites from "../../../api/favorite/getFavorites";

const getAddr = (lat: number, lng: number): Promise<string> => {
  return new Promise((resolve) => {
    let geocoder = new window.kakao.maps.services.Geocoder();
    let coord = new window.kakao.maps.LatLng(lat, lng);

    geocoder.coord2Address(
      coord.getLng(),
      coord.getLat(),
      (
        result: { address: { address_name: string | PromiseLike<string> } }[],
        status: string
      ) => {
        if (status === window.kakao.maps.services.Status.OK) {
          resolve(result[0].address.address_name);
        } else {
          resolve("주소 정보 없음");
        }
      }
    );
  });
};

const useGetFavorites = () => {
  return useQuery({
    queryKey: ["favorite"],
    queryFn: async () => {
      const result = await getFavorites();

      const newMarkers = await Promise.all(
        result.map(async (marker) => ({
          ...marker,
          addr: await getAddr(marker.latitude, marker.longitude),
        }))
      );

      return newMarkers;
    },
    retry: false,
  });
};

export default useGetFavorites;
