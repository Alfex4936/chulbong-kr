import Link from "next/link";

const Unauthenticated = () => {
  return (
    <div>
      {/* <h1 className="text-center mb-4">로그인 해주세요</h1> */}

      <Link
        href={"/signin"}
        className="block w-full text-left group rounded-sm mb-3 px-1 py-2 hover:bg-black-light-2 mo:text-sm"
      >
        <div className="flex justify-center transition-transform duration-75 transform group-hover:scale-95">
          🔑 로그인 하러 가기
        </div>
      </Link>
    </div>
  );
};

export default Unauthenticated;
