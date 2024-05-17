"use client";

import ChangePassword from "@/app/(navbar)/mypage/user/ChangePassword";
import Link from "next/link";

const SigninBottomlinks = () => {
  return (
    <div>
      <div className="flex justify-start items-center mb-1">
        <p className="text-sm mr-1">계정이 없으신가요?</p>
        <Link
          href={"/signup"}
          className="text-sm text-grey-dark hover:underline"
        >
          이메일로 회원가입 하기
        </Link>
      </div>
      <div className="flex justify-start items-center">
        <p className="text-sm mr-1">비밀번호를 잊어버리셨나요?</p>
        <ChangePassword
          textClass="text-sm text-grey-dark hover:underline"
          text="비밀번호 초기화하기"
        />
      </div>
    </div>
  );
};

export default SigninBottomlinks;
