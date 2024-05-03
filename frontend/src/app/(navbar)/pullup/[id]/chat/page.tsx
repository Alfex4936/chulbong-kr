import instance from "@/api/instance";
import BlackSideBody from "@/components/atom/BlackSideBody";
import { type Marker } from "@/types/Marker.types";
import PullupChatClient from "./PullupChatClient";

const getMarker = async (id: number): Promise<Marker> => {
  const res = await instance.get(
    `${process.env.NEXT_PUBLIC_BASE_URL}/markers/${id}/details`
  );

  return res.data;
};

export const generateMetadata = async ({
  params,
}: {
  params: { id: string };
}) => {
  const data = await getMarker(~~params.id);
  return {
    title: `${data.address || "대한민국 철봉 지도"} | 채팅`,
    description: data ? `${data.address} 지역 채팅` : "잘못된 요청입니다.",
  };
};

const PullupChat = ({ params }: { params: { id: string } }) => {
  console.log(params.id);
  return (
    <BlackSideBody toggle bodyClass="p-0 mo:p-0">
      <PullupChatClient markerId={Number(params.id)} />
    </BlackSideBody>
  );
};

export default PullupChat;
