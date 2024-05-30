import instance from "../instance";

const approveReport = async (reportId: number) => {
  const res = await instance.post(`/api/v1/reports/approve/${reportId}`);

  return res.data;
};

export default approveReport;
