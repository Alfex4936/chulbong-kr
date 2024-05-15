import instance from "@/api/instance";
import BlackSideBody from "@/components/atom/BlackSideBody";
import PrevHeader from "@/components/atom/PrevHeader";
import {
  HydrationBoundary,
  QueryClient,
  dehydrate,
} from "@tanstack/react-query";
import { cookies } from "next/headers";
import ReportClient from "./reportClient";
// TODO: 리포트 삭제

const getMyReports = async (cookie: string) => {
  const res = await instance.get(
    `${process.env.NEXT_PUBLIC_BASE_URL}/users/reports`,
    {
      headers: {
        Cookie: cookie || "",
      },
    }
  );

  return res.data;
};

const ReportPage = async () => {
  const cookieStore = cookies();
  const decodeCookie = decodeURIComponent(cookieStore.toString());
  const queryClient = new QueryClient();

  await queryClient.prefetchQuery({
    queryKey: ["marker", "report", "me"],
    queryFn: () => {
      return getMyReports(decodeCookie);
    },
  });

  const dehydrateState = dehydrate(queryClient);
  return (
    <BlackSideBody toggle bodyClass="p-0 mo:px-0 mo:pb-0">
      <PrevHeader url={`/mypage`} text="정보 수정 제안 목록" />
      <HydrationBoundary state={dehydrateState}>
        <div className="px-4 pt-2 mo:pb-20">
          <ReportClient />
        </div>
      </HydrationBoundary>
    </BlackSideBody>
  );
};

export default ReportPage;
