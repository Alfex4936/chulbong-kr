import instance from "@/api/instance";
import { type MyInfo } from "@/api/user/myInfo";
import BlackSideBody from "@/components/atom/BlackSideBody";
import PrevHeader from "@/components/atom/PrevHeader";
import {
  QueryClient,
  dehydrate,
  HydrationBoundary,
} from "@tanstack/react-query";
import { cookies } from "next/headers";
import UserClient from "./UserClient";

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

const User = async () => {
  const cookieStore = cookies();
  const decodeCookie = decodeURIComponent(cookieStore.toString());
  const queryClient = new QueryClient();

  await queryClient.prefetchQuery({
    queryKey: ["user", "me"],
    queryFn: () => {
      return myInfo(decodeCookie);
    },
  });

  const dehydrateState = dehydrate(queryClient);

  return (
    <BlackSideBody>
      <PrevHeader back text="내 정보" />

      <HydrationBoundary state={dehydrateState}>
        <div className="pt-2">
          <UserClient />
        </div>
      </HydrationBoundary>
    </BlackSideBody>
  );
};

export default User;
