import { useQuery } from "@tanstack/react-query";
import getWeather from "../../../api/markers/getWeather";

const useWeatherData = (lat: number, lng: number, start: boolean) => {
  return useQuery({
    queryKey: ["marker", "weather"],
    queryFn: async () => {
      return getWeather(lat, lng);
    },
    retry: false,
    gcTime: 0,
    enabled: start,
  });
};

export default useWeatherData;
