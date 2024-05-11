import instance from "../instance";

export interface ReportsRes {
  createdAt: Date;
  description: string;
  latitude: number;
  longitude: number;
  markerId: number;
  reportImageUrl: string;
  userId: number;
}

const getMyReports = async (): Promise<ReportsRes[]> => {
  const res = await instance.get(`/api/v1/users/reports`);

  return res.data;
};

export default getMyReports;
