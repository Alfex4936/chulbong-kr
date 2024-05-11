import instance from "../instance";
import { type ReportsRes } from "./getMyReports";

const getReportsForMarker = async (markerId: number): Promise<ReportsRes[]> => {
  const res = await instance.get(`/api/v1/reports/marker/${markerId}`);

  return res.data;
};

export default getReportsForMarker;
