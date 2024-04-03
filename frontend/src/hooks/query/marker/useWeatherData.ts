import { useQuery } from "@tanstack/react-query";
import getWeather from "../../../api/markers/getWeather";

const useWeatherData = (
  lat: number,
  lng: number,
  start: boolean,
  id: number
) => {
  return useQuery({
    queryKey: ["marker", "weather", id],
    queryFn: async () => {
      return getWeather(lat, lng);
    },
    retry: false,
    enabled: start,
  });
};

export default useWeatherData;
