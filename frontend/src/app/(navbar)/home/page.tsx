import BlackSideBody from "@/components/atom/BlackSideBody";
import Heading from "@/components/atom/Heading";
import { Separator } from "@/components/ui/separator";
import NoticeSlide from "./_components/NoticeSlide";
import Ranking from "./_components/Ranking";
import SearchInput from "./_components/SearchInput";

const Home = async () => {
  return (
    <BlackSideBody toggle>
      <Heading title="대한민국 철봉 지도" />
      <SearchInput />
      <NoticeSlide />
      <Separator className="my-8 bg-grey-dark-1" />
      <Ranking />
    </BlackSideBody>
  );
};

export default Home;
