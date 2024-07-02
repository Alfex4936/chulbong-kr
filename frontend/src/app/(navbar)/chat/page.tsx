import BlackSideBody from "@/components/atom/BlackSideBody";
import ChatClient from "./ChatClient";

export const generateMetadata = () => {
  return {
    title: "대한민국 철봉 지도 | 채팅",
    description: "지역별 채팅에 참여하세요.",
  };
};

const Chat = () => {
  return (
    <BlackSideBody>
      <ChatClient />
    </BlackSideBody>
  );
};

export default Chat;
