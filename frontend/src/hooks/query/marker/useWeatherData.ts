import { useQuery } from "@tanstack/react-query";
import getWeather from "../../../api/markers/getWeather";

const useWeatherData = (lat: number, lng: number) => {
  return useQuery({
    queryKey: ["marker", "weather"],
    queryFn: async () => {
      return getWeather(lat, lng);
    },
    retry: false,
  });
};

export default useWeatherData;
