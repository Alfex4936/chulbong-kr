import BlackSideBody from "@/components/atom/BlackSideBody";
import Link from "next/link";

const NotFound = () => {
  return (
    <BlackSideBody toggle>
      <div className="flex flex-col items-center">
        <h1 className="text-7xl text-center py-9 mt-5">404</h1>
        <p className="text-center">찾으시는 페이지가 존재하지 않습니다.</p>
        <Link
          href={"/home"}
          className="border border-grey border-solid mt-5 px-3 py-1 rounded-sm hover:bg-white-tp-dark hover:text-black"
        >
          홈으로 가기
        </Link>
      </div>
    </BlackSideBody>
  );
};

export default NotFound;
