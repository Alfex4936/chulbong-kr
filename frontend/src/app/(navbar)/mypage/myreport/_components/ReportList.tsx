import { type Report } from "@/api/report/getReportForMyMarker";
import { ArrowRightIcon } from "@/components/icons/ArrowIcons";
import { Fragment, useState } from "react";
import MarkerReportList from "../../report/_components/MarkerReportList";

interface Props {
  id: number;
  count: number;
  reports: Report[];
  address: string;
  desc: string;
}

const ReportList = ({ id, count, reports, address, desc }: Props) => {
  const [toggle, setToggle] = useState(false);

  return (
    <Fragment>
      <button
        className={`flex items-center justify-between text-left w-full bg-black-light-2 mb-4 p-3 rounded-sm`}
        onClick={() => setToggle((prev) => !prev)}
      >
        <div className="w-full">
          <div className="text-grey-dark">
            총 <span className="font-bold text-grey-light">{count || 0}</span>개
            요청
          </div>
          <div className="break-words">{address}</div>
          <div className="text-grey-dark text-sm break-words">{desc}</div>
        </div>
        <div className={`${toggle ? "rotate-90" : ""} duration-150`}>
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
                imgs={report.photos}
                markerId={id}
                reportId={report.reportID}
                status={report.status}
                address={reports[0].address}
                isAdmin={true}
              />
            );
          })}
        </>
      )}
    </Fragment>
  );
};

export default ReportList;
