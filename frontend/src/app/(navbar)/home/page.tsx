import BlackSideBody from "@/components/atom/BlackSideBody";
import Heading from "@/components/atom/Heading";
import { Separator } from "@/components/ui/separator";
import LinkEmojiButton from "./_components/LinkEmojiButton";
import LocationQuery from "./_components/LocationQuery";
import Ranking from "./_components/Ranking";
import SearchInput from "./_components/SearchInput";

const Home = () => {
  return (
    <BlackSideBody className="px-4">
      <LocationQuery />
      <Heading title="대한민국 철봉 지도" />
      <SearchInput sticky />
      <div className="mt-4">
        <LinkEmojiButton
          url="/search"
          text="주변 검색"
          subText="지도 중앙을 기준으로 주변 위치를 검색하세요!"
          emoji="🔍"
        />
        <LinkEmojiButton
          url="/pullup/register"
          text="위치 등록"
          subText="위치를 등록하고 다른 사람들과 공유하세요!"
          emoji="🚩"
        />
      </div>
      <Separator className="my-8 bg-grey-dark-1" />
      <Ranking />
    </BlackSideBody>
  );
};

export default Home;
