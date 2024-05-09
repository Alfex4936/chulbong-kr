import instance from "../instance";

const reportMarker = async () => {
  const res = await instance.post(`/api/v1/reports`);

  return res.data;
};

export default reportMarker;
