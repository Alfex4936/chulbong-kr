import axios from "axios";

const markerRanking = async () => {
  const res = await axios.get(`/api/v1/markers/ranking`);

  return res.data;
};

export default markerRanking;
