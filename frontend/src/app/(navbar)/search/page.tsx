import BlackSideBody from "@/components/atom/BlackSideBody";
import ErrorMessage from "@/components/atom/ErrorMessage";
import Heading from "@/components/atom/Heading";

const Search = () => {
  return (
    <BlackSideBody toggle>
      <Heading title="주변 검색" />
      <ErrorMessage
        text="거리는 부정확할 수 있고, 현재 보이는 화면 중아에서부터 찾습니다."
        className="text-center"
      />
    </BlackSideBody>
  );
};

export default Search;
