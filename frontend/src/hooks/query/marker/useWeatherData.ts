import getWeather from "@/api/markers/getWeather";
import { useQuery } from "@tanstack/react-query";

const useWeatherData = (lat: number, lng: number, start: boolean) => {
  return useQuery({
    queryKey: ["weather", lat, lng],
    queryFn: () => {
      return getWeather(lat, lng);
    },

    enabled: start,
  });
};

export default useWeatherData;
