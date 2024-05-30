import instance from "../instance";

const denyReport = async (reportId: number) => {
  const res = await instance.post(`/api/v1/reports/deny/${reportId}`);

  return res.data;
};

export default denyReport;
