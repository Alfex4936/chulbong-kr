import instance from "@/api/instance";
import { type MyInfo } from "@/api/user/myInfo";
import BlackSideBody from "@/components/atom/BlackSideBody";
import Heading from "@/components/atom/Heading";
import {
  HydrationBoundary,
  QueryClient,
  dehydrate,
} from "@tanstack/react-query";
import { cookies } from "next/headers";
import MypageClient from "./MypageClient";
// TODO: 내정보 캐싱 확인

const myInfo = async (cookie: string): Promise<MyInfo> => {
  const res = await instance.get(
    `${process.env.NEXT_PUBLIC_BASE_URL}/users/me`,
    {
      headers: {
        Cookie: cookie || "",
      },
    }
  );

  return res.data;
};
// const checkAdmin = async (cookie: string) => {
//   const res = await instance.get(`${process.env.NEXT_PUBLIC_BASE_URL}/admin`, {
//     headers: {
//       Cookie: cookie || "",
//     },
//   });

//   return res.data;
// };

export const generateMetadata = () => {
  return {
    title: `대한민국 철봉 지도 | 마이 페이지`,
  };
};

const Mypage = async () => {
  const cookieStore = cookies();
  const decodeCookie = decodeURIComponent(cookieStore.toString());
  const queryClient = new QueryClient();

  await queryClient.invalidateQueries({ queryKey: ["user", "me"] });

  await queryClient.prefetchQuery({
    queryKey: ["user", "me"],
    queryFn: () => {
      return myInfo(decodeCookie);
    },
  });
  // await queryClient.prefetchQuery({
  //   queryKey: ["admin", "check"],
  //   queryFn: () => {
  //     return checkAdmin(decodeCookie);
  //   },
  // });

  const dehydrateState = dehydrate(queryClient);

  return (
    <BlackSideBody toggle>
      <Heading title="내 정보" />
      <HydrationBoundary state={dehydrateState}>
        <MypageClient />
      </HydrationBoundary>
    </BlackSideBody>
  );
};

export default Mypage;
