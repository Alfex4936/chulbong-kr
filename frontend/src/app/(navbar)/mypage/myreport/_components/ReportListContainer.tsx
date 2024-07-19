import type { MyMarkerReportRes } from "@/api/report/getReportForMyMarker";
import ReportList from "./ReportList";

interface Props {
  data: MyMarkerReportRes;
}

const ReportListContainer = ({ data }: Props) => {
  // TODO: 여기 주소랑 설명 data에 가끔 안보임
  const reportItems = Object.entries(data.markers).map(([key, reports]) => (
    <div key={key}>
      <ReportList
        id={Number(key)}
        count={reports.length}
        reports={reports}
        address={reports[0].address || "주소 제공 안됨"}
        desc={`위치 id: ${key}` || "작성 된 설명이 없습니다"}
      />
    </div>
  ));
  return <div>{reportItems}</div>;
};

export default ReportListContainer;
