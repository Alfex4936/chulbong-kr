import { useInfiniteQuery } from "@tanstack/react-query";
import getMyMarker from "../../api/markers/getMyMarker";

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

const useGetMyMarker = () => {
  return useInfiniteQuery({
    queryKey: ["myMarker"],
    queryFn: async ({ pageParam }) => {
      const result = await getMyMarker({ pageParam });

      if (result.markers) {
        const newMarkers = await Promise.all(
          result.markers.map(async (marker) => ({
            ...marker,
            addr: await getAddr(marker.latitude, marker.longitude),
          }))
        );

        result.markers = newMarkers;
      }

      return result;
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage) => {
      if (lastPage.currentPage < lastPage.totalPages)
        return lastPage.currentPage + 1;
      return undefined;
    },
    retry: false,
  });
};

export default useGetMyMarker;
