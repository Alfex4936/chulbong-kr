import instance from "../instance";

const approveReport = async (markerId: number) => {
  const res = await instance.post(`/api/v1/reports/approve/${markerId}`);

  return res.data;
};

export default approveReport;
