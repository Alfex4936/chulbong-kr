import instance from "@/api/instance";
import Tabs from "@/components/atom/Tabs";
import {
  HydrationBoundary,
  QueryClient,
  dehydrate,
} from "@tanstack/react-query";
import RankingResult from "./RankingResult";

const Ranking = async () => {
  const queryClient = new QueryClient();

  const fetchRanking = async () => {
    const res = await instance.get(
      `${process.env.NEXT_PUBLIC_BASE_URL}/markers/ranking`
    );

    return res.data;
  };

  await queryClient.prefetchQuery({
    queryKey: ["ranking", "top10"],
    queryFn: fetchRanking,
  });

  const dehydrateState = dehydrate(queryClient);

  return (
    <div>
      <HydrationBoundary state={dehydrateState}>
        <Tabs title="랭킹" tabs={["전체", "주변"]}>
          <RankingResult />
        </Tabs>
      </HydrationBoundary>
    </div>
  );
};

export default Ranking;
