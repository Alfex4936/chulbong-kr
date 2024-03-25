import axios from "axios";

export interface RankingInfo {
  address: string;
  latitude: number;
  longitude: number;
  makerId: number;
}

const markerRanking = async (): Promise<RankingInfo[]> => {
  const res = await axios.get(`/api/v1/markers/ranking`);

  return res.data;
};

export default markerRanking;
