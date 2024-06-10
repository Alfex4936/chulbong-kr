import { type Report } from "@/api/report/getReportForMyMarker";
import { ArrowRightIcon } from "@/components/icons/ArrowIcons";
import { Skeleton } from "@/components/ui/skeleton";
import useMarkerData from "@/hooks/query/marker/useMarkerData";
import { Fragment, useState } from "react";
import MarkerReportList from "../../report/_components/MarkerReportList";

interface Props {
  id: number;
  count: number;
  reports: Report[];
}

const ReportList = ({ id, count, reports }: Props) => {
  const { data, isFetching } = useMarkerData(id);

  const [toggle, setToggle] = useState(false);

  if (!data)
    return <Skeleton className="bg-black-light-2 mb-4 p-2 rounded-sm" />;

  return (
    <Fragment>
      <button
        className="flex items-center justify-between text-left w-full bg-black-light-2 mb-4 p-3 rounded-sm"
        onClick={() => setToggle((prev) => !prev)}
      >
        <div>
          <div className="text-grey-dark">
            총 <span className="font-bold text-grey-light">{count || 0}</span>개
            요청 중
          </div>
          <div className="truncate">{data.address}</div>
          <div className="truncate text-grey-dark text-sm">
            {data.description || "작성 된 설명이 없습니다."}
          </div>
        </div>
        <div>
          <ArrowRightIcon />
        </div>
      </button>

      {toggle && (
        <>
          {reports.map((report) => {
            return (
              <MarkerReportList
                className="m-0 mb-4 w-full"
                key={report.reportID}
                desc={report.description}
                img={report.photos}
                lat={data.latitude}
                lng={data.longitude}
                markerId={id}
                reportId={report.reportID}
                status={report.status}
                isFetching={isFetching}
              />
            );
          })}
        </>
      )}
    </Fragment>
  );
};

export default ReportList;
