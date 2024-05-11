"use client";

import useGetMyReports from "@/hooks/query/report/useGetMyReports";
import { isAxiosError } from "axios";
import MarkerReportList from "./_components/MarkerReportList";

const ReportClient = () => {
  const { data: myReports, error, isError } = useGetMyReports();

  if (isError) {
    if (isAxiosError(error)) {
      if (error.response?.status === 404) {
        return (
          <div className="text-center">정보 수정 제안한 위치가 없습니다.</div>
        );
      } else {
        return <div className="text-center">잠시 후 다시 시도해 주세요.</div>;
      }
    } else {
      return <div className="text-center">잠시 후 다시 시도해 주세요.</div>;
    }
  }
  // TODO: 이미지 개수 이상 오류
  // console.log(myReports);

  return (
    <div>
      {myReports?.map((report) => {
        return (
          <div key={report.reportId} className="mb-4">
            <MarkerReportList
              markerId={report.markerId}
              lat={report.latitude}
              lng={report.longitude}
              desc={report.description}
              img={report.photoUrls}
              status={report.status}
            />
          </div>
        );
      })}
    </div>
  );
};

export default ReportClient;
