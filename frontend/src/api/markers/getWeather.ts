import instance from "../instance";

export interface WeatherRes {
  temperature: string;
  desc: string;
  humidity: string;
  rainfall: string;
  snowfall: string;
  iconImage: string;
}

const getWeather = async (lat: number, lng: number): Promise<WeatherRes> => {
  const res = await instance.get(
    `/api/v1/markers/weather?latitude=${lat}&longitude=${lng}`
  );

  return res.data;
};

export default getWeather;
