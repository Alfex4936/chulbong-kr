import Image from "next/image";
import Link from "next/link";
import { Suspense } from "react";
import SigninBottomlinks from "./_components/SigninBottomlinks";
import SigninForm from "./_components/signin-form";

const Signin = () => {
  return (
    <div>
      <div className="flex items-center justify-center mt-14 w-full">
        <Image src={"/2.png"} alt="logo" width={45} height={45} className="mr-2" />
        <Link href={"/home"} className="text-2xl">
          대한민국 철봉 지도
        </Link>
      </div>
      <div className="mx-auto w-full max-w-[500px] min-w-80 p-10">
        <Suspense>
          <SigninForm />
        </Suspense>
      </div>
      <div className="mx-auto w-full max-w-[500px] min-w-80 px-10">
        <SigninBottomlinks />
      </div>
    </div>
  );
};

export default Signin;
