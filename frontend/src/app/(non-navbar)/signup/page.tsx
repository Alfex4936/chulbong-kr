import Link from "next/link";
import SignupForm from "./_components/signup-form";

const Signup = () => {
  return (
    <div>
      <Link
        href={"/home"}
        className="inline-block w-full text-center text-2xl mt-14"
      >
        대한민국 철봉 지도
      </Link>
      <div className="mx-auto w-full max-w-[500px] min-w-80 p-10">
        <SignupForm />
      </div>
    </div>
  );
};

export default Signup;
