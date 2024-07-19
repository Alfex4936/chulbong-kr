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
    <BlackSideBody>
      <PrevHeader back text="정보 수정 제안 목록" />
      <HydrationBoundary state={dehydrateState}>
        <div className="pt-2">
          <ReportClient />
        </div>
      </HydrationBoundary>
    </BlackSideBody>
  );
};

export default ReportPage;
