"use client";

import useGetMyReports from "@/hooks/query/report/useGetMyReports";
import { isAxiosError } from "axios";

const ReportClient = () => {
  const { data: myReports, error, isError } = useGetMyReports();

  console.log(myReports);

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
  return <div>reportClient</div>;
};

export default ReportClient;
