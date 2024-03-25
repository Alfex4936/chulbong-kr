import axios from "axios";

const markerAriaRanking = async (lat: number, lng: number) => {
  const res = await axios.get(
    `/api/v1/markers/area-ranking?latitude=${lat}&longitude=${lng}&limit=10`
  );

  return res.data;
};

export default markerAriaRanking;
