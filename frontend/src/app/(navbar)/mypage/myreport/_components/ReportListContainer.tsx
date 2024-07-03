import type { MyMarkerReportRes } from "@/api/report/getReportForMyMarker";
import ReportList from "./ReportList";

interface Props {
  data: MyMarkerReportRes;
}

const ReportListContainer = ({ data }: Props) => {
  const reportItems = Object.entries(data.markers).map(([key, reports]) => (
    <div key={key}>
      <ReportList
        id={Number(key)}
        count={reports.length}
        reports={reports}
        address={reports[0].address || "주소 제공 안됨"}
        desc={reports[0].description || "작성 된 설명이 없습니다"}
      />
    </div>
  ));
  return <div>{reportItems}</div>;
};

export default ReportListContainer;
