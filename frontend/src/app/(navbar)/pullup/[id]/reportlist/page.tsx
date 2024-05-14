import instance from "@/api/instance";
import ReportClient from "@/app/(navbar)/mypage/report/reportClient";
import BlackSideBody from "@/components/atom/BlackSideBody";
import PrevHeader from "@/components/atom/PrevHeader";
import {
  HydrationBoundary,
  QueryClient,
  dehydrate,
} from "@tanstack/react-query";
import Link from "next/link";
import LinkWrap from "./LinkWrap";
// TODO: 프리페치 안되는 문제 해결
// TODO: 등록한 유저면 승인 or 거절
// TODO: 승인 거절 삭제 안됨

const getMarkerReport = async (markerId: number) => {
  const res = await instance.get(
    `${process.env.NEXT_PUBLIC_BASE_URL}/reports/marker/${markerId}`
  );

  return res.data;
};

const ReportListPage = async ({ params }: { params: { id: string } }) => {
  const queryClient = new QueryClient();

  await queryClient.prefetchQuery({
    queryKey: ["marker", "report", "formarker", params.id],
    queryFn: () => {
      return getMarkerReport(Number(params.id));
    },
  });

  const dehydrateState = dehydrate(queryClient);
  return (
    <BlackSideBody toggle bodyClass="relative p-0 mo:px-0 mo:pb-0">
      <PrevHeader url={`/pullup/${params.id}`} text="정보 수정 제안 목록" />

      <div className="px-4 pb-4 scrollbar-thin mo:pb-20">
        <LinkWrap id={params.id} />
        <HydrationBoundary state={dehydrateState}>
          <ReportClient type="formarker" markerId={Number(params.id)} />
        </HydrationBoundary>
      </div>
    </BlackSideBody>
  );
};

export default ReportListPage;
