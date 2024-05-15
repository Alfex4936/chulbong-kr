import instance from "@/api/instance";
import { type Favorite } from "@/api/user/bookmarkMarker";
import BlackSideBody from "@/components/atom/BlackSideBody";
import PrevHeader from "@/components/atom/PrevHeader";
import {
  HydrationBoundary,
  QueryClient,
  dehydrate,
} from "@tanstack/react-query";
import { cookies } from "next/headers";
import BookmarkClient from "./BookmarkClient";

const bookmarkMarker = async (cookie: string): Promise<Favorite> => {
  const res = await instance.get(
    `${process.env.NEXT_PUBLIC_BASE_URL}/users/favorites`,
    {
      headers: {
        Cookie: cookie || "",
      },
    }
  );

  return res.data;
};

const Bookmark = async () => {
  const cookieStore = cookies();
  const decodeCookie = decodeURIComponent(cookieStore.toString());
  const queryClient = new QueryClient();

  await queryClient.prefetchQuery({
    queryKey: ["marker", "bookmark"],
    queryFn: () => {
      return bookmarkMarker(decodeCookie);
    },
  });

  const dehydrateState = dehydrate(queryClient);

  return (
    <BlackSideBody toggle bodyClass="p-0 mo:px-0 mo:pb-0">
      <PrevHeader url="/mypage" text="저장한 장소" />

      <HydrationBoundary state={dehydrateState}>
        <div className="px-4 pt-2 pb-4 mo:pb-20">
          <BookmarkClient />
        </div>
      </HydrationBoundary>
    </BlackSideBody>
  );
};

export default Bookmark;
