import instance from "../instance";

export interface RankingInfo {
  address: string;
  latitude: number;
  longitude: number;
  markerId: number;
}

const getMarkerRanking = async (): Promise<RankingInfo[]> => {
  const res = await instance.get(`/api/v1/markers/ranking`);

  return res.data;
};

export default getMarkerRanking;
