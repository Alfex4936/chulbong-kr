import BlackSideBody from "@/components/atom/BlackSideBody";
import GrowBox from "@/components/atom/GrowBox";
import Heading from "@/components/atom/Heading";
import { Separator } from "@/components/ui/separator";
import Link from "next/link";
import NoticeSlide from "./_components/NoticeSlide";
import Ranking from "./_components/Ranking";
import SearchInput from "./_components/SearchInput";

const Home = async () => {
  return (
    <BlackSideBody toggle bodyClass="px-6">
      <Heading title="대한민국 철봉 지도" />
      <SearchInput />
      <NoticeSlide />
      <div className="mt-4">
        <Link href={"/pullup/regist"}>
          <div className="block w-full text-left group rounded-sm mb-3 px-1 py-2 hover:bg-black-light-2 text-sm">
            <div
              className={`flex justify-start transition-transform duration-75 transform group-hover:scale-95`}
            >
              <span className="mr-2">🚩</span>

              <span>위치 등록</span>
              <GrowBox />
              <span className="text-grey-dark-1 text-xs">
                위치를 등록하고 다른 사람들과 공유하세요!
              </span>
            </div>
          </div>
        </Link>
      </div>
      <Separator className="my-8 bg-grey-dark-1" />
      <Ranking />
    </BlackSideBody>
  );
};

export default Home;
