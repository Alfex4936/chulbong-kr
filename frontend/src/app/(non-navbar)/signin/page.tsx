"use client";

import SigninForm from "./_components/signin-form";

const Signin = () => {
  return (
    <div>
      <h1 className="absolute top-24 left-1/2 -translate-x-1/2 text-center text-2xl">
        대한민국 철봉 지도
      </h1>
      <div className="absolute top-40 left-1/2 -translate-x-1/2 w-full max-w-[500px] min-w-80 p-10">
        <SigninForm />
      </div>
    </div>
  );
};

export default Signin;
