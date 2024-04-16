import BlackSideBody from "@/components/atom/BlackSideBody";
import Heading from "@/components/atom/Heading";
import NoticeSlide from "./_components/NoticeSlide";
import SearchInput from "./_components/SearchInput";

const Home = () => {
  return (
    <BlackSideBody toggle>
      <Heading title="대한민국 철봉 지도" />
      <NoticeSlide />
      <SearchInput />
    </BlackSideBody>
  );
};

export default Home;
