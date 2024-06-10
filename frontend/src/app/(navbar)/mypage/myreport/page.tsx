import instance from "@/api/instance";
import BlackSideBody from "@/components/atom/BlackSideBody";
import PrevHeader from "@/components/atom/PrevHeader";
import {
    HydrationBoundary,
    QueryClient,
    dehydrate,
} from "@tanstack/react-query";
import { cookies } from "next/headers";
import MyreportClient from "./myreportClient";

const getReportForMyMarker = async (cookie: string) => {
  const res = await instance.get(
    `${process.env.NEXT_PUBLIC_BASE_URL}/users/reports/for-my-markers`,
    {
      headers: {
        Cookie: cookie || "",
      },
    }
  );

  return res.data;
};

const MyReportPage = async () => {
  const cookieStore = cookies();
  const decodeCookie = decodeURIComponent(cookieStore.toString());
  const queryClient = new QueryClient();

  await queryClient.prefetchQuery({
    queryKey: ["report", "my"],
    queryFn: () => {
      return getReportForMyMarker(decodeCookie);
    },
  });

  const dehydrateState = dehydrate(queryClient);

  return (
    <BlackSideBody toggle bodyClass="p-0 mo:px-0 mo:pb-0">
      <PrevHeader url={`/mypage`} text="정보 수정 제안 목록" />
      <HydrationBoundary state={dehydrateState}>
        <div className="px-4 pt-2 pb-4 mo:pb-20">
          <MyreportClient />
        </div>
      </HydrationBoundary>
    </BlackSideBody>
  );
};

export default MyReportPage;
