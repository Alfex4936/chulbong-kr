import instance from "@/api/instance";
import { type MyMarkerRes } from "@/api/user/mylocateMarker";
import BlackSideBody from "@/components/atom/BlackSideBody";
import Heading from "@/components/atom/Heading";
import PrevHeader from "@/components/atom/PrevHeader";
import {
  HydrationBoundary,
  QueryClient,
  dehydrate,
} from "@tanstack/react-query";
import { cookies } from "next/headers";
import MylocateClient from "./MylocateClient";

const mylocateMarker = async (
  pageParam: number,
  cookie: string
): Promise<MyMarkerRes> => {
  const res = await instance.get(
    `${process.env.NEXT_PUBLIC_BASE_URL}/markers/my?page=${pageParam}`,
    {
      headers: {
        Cookie: cookie || "",
      },
    }
  );

  return res.data;
};

const Mylocate = async () => {
  const cookieStore = cookies();
  const decodeCookie = decodeURIComponent(cookieStore.toString());
  const queryClient = new QueryClient();

  await queryClient.prefetchInfiniteQuery({
    queryKey: ["myMarker"],
    queryFn: ({ pageParam = 1 }) => mylocateMarker(pageParam, decodeCookie),
    initialPageParam: 1,
  });

  const dehydrateState = dehydrate(queryClient);

  return (
    <BlackSideBody toggle bodyClass="p-0 mo:px-0 mo:pb-0">
      <PrevHeader url="/mypage" text="내 정보" />

      <Heading title="등록한 장소" />
      <HydrationBoundary state={dehydrateState}>
        <MylocateClient />
      </HydrationBoundary>
    </BlackSideBody>
  );
};

export default Mylocate;
