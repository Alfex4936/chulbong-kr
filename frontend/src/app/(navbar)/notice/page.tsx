import BlackSideBody from "@/components/atom/BlackSideBody";
import Heading from "@/components/atom/Heading";
import NoticeClient from "./NoticeClient";

const Notice = () => {
  return (
    <BlackSideBody>
      <Heading title="공지" />
      <NoticeClient />
    </BlackSideBody>
  );
};

export default Notice;
