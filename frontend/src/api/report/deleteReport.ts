import instance from "../instance";

const deleteReport = async (markerId: number, reportId: number) => {
  const res = await instance.delete(
    `/api/v1/reports?markerID=${markerId}&reportID=${reportId}`
  );

  return res.data;
};

export default deleteReport;
