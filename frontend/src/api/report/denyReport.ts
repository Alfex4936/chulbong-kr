import instance from "../instance";

const denyReport = async (markerId: number) => {
  const res = await instance.post(`/api/v1/reports/deny/${markerId}`);

  return res.data;
};

export default denyReport;
