"use client";

import { type ReportsRes } from "@/api/report/getMyReports";
import useReportsData from "@/hooks/query/report/useReportsData";
import useMyinfoData from "@/hooks/query/user/useMyinfoData";
import { QueryObserverRefetchErrorResult } from "@tanstack/react-query";
import { isAxiosError } from "axios";
import MarkerReportList from "./_components/MarkerReportList";

interface Props {
  type?: "me" | "formarker" | "all";
  markerId?: number;
}

const ReportClient = ({ type = "me", markerId }: Props) => {
  const {
    data: myReports,
    error,
    isError,
  } = useReportsData({ type, markerId }) as QueryObserverRefetchErrorResult<
    ReportsRes[],
    Error
  >;

  console.log(myReports);

  const { data: myInfo } = useMyinfoData();

  console.log(myInfo);

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

  if (myReports?.length === 0)
    return <div className="text-center">요청중인 제안이 없습니다.</div>;

  return (
    <div>
      {myReports?.map((report) => {
        return (
          <div key={report.reportId} className="mb-4">
            <MarkerReportList
              markerId={report.markerId}
              desc={report.description}
              imgs={report.photoUrls}
              status={report.status}
              userId={report.userId}
              reportId={report.reportId}
              myId={myInfo?.userId}
              address={report.address}
              isAdmin={report.userId === myInfo?.userId || myInfo?.chulbong}
            />
          </div>
        );
      })}
    </div>
  );
};

export default ReportClient;
