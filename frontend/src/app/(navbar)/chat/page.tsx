import BlackSideBody from "@/components/atom/BlackSideBody";
import Heading from "@/components/atom/Heading";
import ChatClient from "./ChatClient";

interface Query {
  ci: number;
}

interface Props {
  searchParams: Query;
}

export const generateMetadata = () => {
  return {
    title: "대한민국 철봉 지도 | 채팅",
    description: "지역별 채팅에 참여하세요.",
  };
};

const Chat = ({ searchParams }: Props) => {
  const { ci } = searchParams;

  return (
    <BlackSideBody toggle bodyClass="p-0 mo:p-0">
      <ChatClient />
    </BlackSideBody>
  );
};

export default Chat;
