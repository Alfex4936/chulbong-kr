import instance from "../instance";
import { type ReportsRes } from "./getMyReports";

const getAllReports = async (): Promise<ReportsRes[]> => {
  const res = await instance.get(`/api/v1/reports/all`);

  return res.data;
};

export default getAllReports;
