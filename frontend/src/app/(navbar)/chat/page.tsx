import BlackSideBody from "@/components/atom/BlackSideBody";
import Heading from "@/components/atom/Heading";
// TODO: 페이지별 메타 태그 연결

interface Query {
  ci: number;
}

interface Props {
  searchParams: Query;
}

const Chat = ({ searchParams }: Props) => {
  const { ci } = searchParams;

  return (
    <BlackSideBody toggle>
      <Heading title={`${ci || "서울"} 채팅방`} subTitle="1명 접속 중" />
    </BlackSideBody>
  );
};

export default Chat;
