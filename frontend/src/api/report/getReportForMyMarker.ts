import instance from "../instance";

export interface Report {
  reportID: number;
  description: string;
  status: string;
  createdAt: string;
  photos: string[];
  address: string;
}

interface Marker {
  [key: string]: Report[];
}

export interface MyMarkerReportRes {
  totalReports: number;
  markers: Marker;
}
// TODO: 여기에 뉴 위치 경도 위도 받기 (승인 시 이동시키기 위함)
const getReportForMyMarker = async (): Promise<MyMarkerReportRes> => {
  const res = await instance.get(`/api/v1/users/reports/for-my-markers`);

  return res.data;
};

export default getReportForMyMarker;
