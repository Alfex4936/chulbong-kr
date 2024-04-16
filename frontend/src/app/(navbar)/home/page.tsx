import BlackSideBody from "@/components/atom/BlackSideBody";
import Heading from "@/components/atom/Heading";
import NoticeSlide from "./_components/NoticeSlide";
import SearchInput from "./_components/SearchInput";
import { Separator } from "@/components/ui/separator";
import Ranking from "./_components/Ranking";

const Home = () => {
  return (
    <BlackSideBody toggle>
      <Heading title="대한민국 철봉 지도" />
      <SearchInput />
      <NoticeSlide />
      <Separator className="my-8 bg-grey-dark" />
      <Ranking />
    </BlackSideBody>
  );
};

export default Home;
