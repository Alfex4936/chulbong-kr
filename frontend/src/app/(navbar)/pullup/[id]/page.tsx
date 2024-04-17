import instance from "@/api/instance";
import BlackSideBody from "@/components/atom/BlackSideBody";
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

// const getMarker = async (id: number): Promise<Marker> => {
//   const res = await instance.get(
//     `${process.env.NEXT_PUBLIC_BASE_URL}/markers/${id}/details`
//   );

//   return res.data;
// };

// export const generateMetadata = async ({ params }: Props) => {
//   const { id } = params;

//   try {
//     const { address, description, favCount } = await getMarker(Number(id));

//     return {
//       title: `철봉 | ${address}`,
//       description: `즐거운 맨몸운동 생활 - ${description} - ${address} - 좋아요 : ${favCount}`,
//     };
//   } catch (error) {
//     return {
//       title: `등록된 위치가 없습니다....`,
//     };
//   }
// };

const Pullup = async ({ params }: Props) => {
  const { id } = params;

  // const queryClient = new QueryClient();

  // queryClient.prefetchQuery({
  //   queryKey: ["marker", id],
  //   queryFn: () => {
  //     return getMarker(Number(id));
  //   },
  // });

  // const dehydrateState = dehydrate(queryClient);

  return (
    <BlackSideBody bodyClass="p-0 mo:px-0 mo:pb-0" toggle>
      {/* <HydrationBoundary state={dehydrateState}> */}
      <PullupClient markerId={Number(id)} />
      {/* </HydrationBoundary> */}
    </BlackSideBody>
  );
};

export default Pullup;
