import Tabs from "@/components/atom/Tabs";
import RankingResult from "./RankingResult";

type Props = {};

const Ranking = () => {


  return (
    <div>
      <Tabs title="랭킹" tabs={["TOP10", "주변"]}>
        <RankingResult />
      </Tabs>
    </div>
  );
};

export default Ranking;
