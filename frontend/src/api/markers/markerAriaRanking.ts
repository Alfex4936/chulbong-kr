import instance from "../instance";
import type { RankingInfo } from "./markerRanking";

const markerAriaRanking = async (
  lat: number,
  lng: number
): Promise<RankingInfo[]> => {
  const res = await instance.get(
    `/api/v1/markers/area-ranking?latitude=${lat}&longitude=${lng}&limit=10`
  );

  return res.data;
};

export default markerAriaRanking;
