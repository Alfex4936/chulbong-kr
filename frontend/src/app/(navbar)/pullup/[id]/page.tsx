import { type CommentsRes } from "@/api/comments/getComments";
import instance from "@/api/instance";
import { type FacilitiesRes } from "@/api/markers/getFacilities";
import BlackSideBody from "@/components/atom/BlackSideBody";
import PrevHeader from "@/components/atom/PrevHeader";
import { type Marker } from "@/types/Marker.types";
import {
  HydrationBoundary,
  QueryClient,
  dehydrate,
} from "@tanstack/react-query";
import PullupClient from "./pullupClient";

type Params = {
  id: string;
};

interface Props {
  params: Params;
}

const getMarker = async (id: number): Promise<Marker> => {
  const res = await instance.get(
    `${process.env.NEXT_PUBLIC_BASE_URL}/markers/${id}/details`
  );

  return res.data;
};

const getComments = async (
  id: number,
  pageParam: number
): Promise<CommentsRes> => {
  const res = await instance.get(
    `${process.env.NEXT_PUBLIC_BASE_URL}/comments/${id}/comments?page=${pageParam}`
  );

  return res.data;
};

const getFacilities = async (markerId: number): Promise<FacilitiesRes[]> => {
  const res = await instance.get(
    `${process.env.NEXT_PUBLIC_BASE_URL}/markers/${markerId}/facilities`
  );

  return res.data;
};

export const generateMetadata = async ({ params }: Props) => {
  const { id } = params;

  try {
    const { address, description, favCount, photos } = await getMarker(
      Number(id)
    );

    return {
      title: `${address} | 철봉`,
      description: `즐거운 맨몸운동 생활 - ${description} - ${address} - 좋아요 : ${favCount}`,
      keywords: `철봉, ${address}`,
      openGraph: {
        type: "website",
        url: `https://www.k-pullup.com/pullup/${id}`,
        title: `${address} | 철봉`,
        description: `즐거운 맨몸운동 생활 - ${description} - ${address} - 좋아요 : ${favCount}`,
        images: photos ? photos[0].photoUrl : "/images/metaimg.webp",
      },
      twitter: {
        card: "summary_large_image",
        title: `${address} | 철봉`,
        description: `즐거운 맨몸운동 생활 - ${description} - ${address} - 좋아요 : ${favCount}`,
        images: photos ? photos[0].photoUrl : "/images/metaimg.webp",
      },
    };
  } catch (error) {
    return {
      title: `등록된 위치가 없습니다....`,
    };
  }
};

const Pullup = async ({ params }: Props) => {
  const { id } = params;

  const queryClient = new QueryClient();

  await queryClient.prefetchQuery({
    queryKey: ["marker", Number(id)],
    queryFn: () => {
      return getMarker(Number(id));
    },
  });

  await queryClient.prefetchQuery({
    queryKey: ["facilities", Number(id)],
    queryFn: () => {
      return getFacilities(Number(id));
    },
  });

  await queryClient.prefetchInfiniteQuery({
    queryKey: ["comments", id],
    queryFn: ({ pageParam = 1 }) => getComments(Number(params.id), pageParam),
    initialPageParam: 1,
  });

  const dehydrateState = dehydrate(queryClient);

  return (
    <BlackSideBody className="flex flex-col">
      <PrevHeader text="상세보기" back />
      <HydrationBoundary state={dehydrateState}>
        <PullupClient markerId={Number(id)} />
      </HydrationBoundary>
    </BlackSideBody>
  );
};

export default Pullup;
