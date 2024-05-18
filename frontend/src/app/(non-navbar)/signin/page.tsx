import Link from "next/link";
import SigninBottomlinks from "./_components/SigninBottomlinks";
import SigninForm from "./_components/signin-form";

const Signin = () => {
  return (
    <div>
      <Link
        href={"/home"}
        className="inline-block w-full text-center text-2xl mt-14"
      >
        대한민국 철봉 지도
      </Link>
      <div className="mx-auto w-full max-w-[500px] min-w-80 p-10">
        <SigninForm />
      </div>
      <div className="mx-auto w-full max-w-[500px] min-w-80 px-10">
        <SigninBottomlinks />
      </div>
    </div>
  );
};

export default Signin;
