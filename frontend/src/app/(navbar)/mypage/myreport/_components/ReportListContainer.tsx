import { type MyMarkerReportRes } from "@/api/report/getReportForMyMarker";
import ReportList from "./ReportList";

interface Props {
  data: MyMarkerReportRes;
}

const ReportListContainer = ({ data }: Props) => {
  const reportItems = Object.entries(data.markers).map(([key, reports]) => (
    <div key={key}>
      <ReportList id={Number(key)} count={reports.length} reports={reports} />
    </div>
  ));
  return <div>{reportItems}</div>;
};

export default ReportListContainer;
