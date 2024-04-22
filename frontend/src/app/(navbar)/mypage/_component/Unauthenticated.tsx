import Link from "next/link";

const Unauthenticated = () => {
  return (
    <div>
      <h1>로그인 해주세요</h1>
      <Link href={"/signin"}>로그인 하러 가기</Link>
    </div>
  );
};

export default Unauthenticated;
