import BlackSideBody from "@/components/atom/BlackSideBody";
import ErrorMessage from "@/components/atom/ErrorMessage";
import PrevHeader from "@/components/atom/PrevHeader";
import SearchRangebar from "./_components/SearchRangebar";

const Search = () => {
  return (
    <BlackSideBody>
      <PrevHeader url="/home" text="주변 검색" />
      <div className="px-4 pt-2 pb-4">
        <ErrorMessage
          text="거리는 부정확할 수 있고, 현재 보이는 화면 중앙에서부터 찾습니다."
          className="text-center"
        />
        <SearchRangebar />
      </div>
    </BlackSideBody>
  );
};

export default Search;
