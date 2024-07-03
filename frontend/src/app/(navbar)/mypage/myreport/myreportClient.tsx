"use client";

import useMyMarkerReportData from "@/hooks/query/report/useMyMarkerReportData";
import React from "react";
import ReportListContainer from "./_components/ReportListContainer";
import { isAxiosError } from "axios";
import { Skeleton } from "@/components/ui/skeleton";

const MyreportClient = () => {
  const { data, isError, error, isLoading } = useMyMarkerReportData();
  
  if (isError) {
    if (isAxiosError(error)) {
      if (error.response?.status === 404) {
        return <div className="text-center">받은 요청이 없습니다.</div>;
      } else {
        return <div className="text-center">잠시 후 다시 시도해 주세요.</div>;
      }
    } else {
      return <div className="text-center">잠시 후 다시 시도해 주세요.</div>;
    }
  }

  if (isLoading)
    <Skeleton className="bg-black-light-2 mb-4 w-full h-20 rounded-sm" />;

  if (!data) return <div className="text-center">받은 요청이 없습니다.</div>;

  return (
    <div>
      <ReportListContainer data={data} />
    </div>
  );
};

export default MyreportClient;
