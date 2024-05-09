import instance from "@/api/instance";
import BlackSideBody from "@/components/atom/BlackSideBody";
import Heading from "@/components/atom/Heading";
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
    queryKey: ["marker", "report"],
    queryFn: () => {
      return getMyReports(decodeCookie);
    },
  });

  const dehydrateState = dehydrate(queryClient);
  return (
    <BlackSideBody toggle>
      <Heading title="정보 수정 제안 목록" />
      <HydrationBoundary state={dehydrateState}>
        <ReportClient />
      </HydrationBoundary>
    </BlackSideBody>
  );
};

export default ReportPage;
