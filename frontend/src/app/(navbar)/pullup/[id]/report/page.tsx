import BlackSideBody from "@/components/atom/BlackSideBody";
import Heading from "@/components/atom/Heading";
import ReportClient from "./reportClient";

const RportMarkerPage = () => {
  return (
    <BlackSideBody toggle>
      <Heading title="정보 수정 제안" />
      <ReportClient />
    </BlackSideBody>
  );
};

export default RportMarkerPage;
