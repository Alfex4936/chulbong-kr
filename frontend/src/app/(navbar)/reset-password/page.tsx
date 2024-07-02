import BlackSideBody from "@/components/atom/BlackSideBody";
import Heading from "@/components/atom/Heading";
import ResetPasswordClient from "./ResetPasswordClient";

const ResetPassword = () => {
  return (
    <BlackSideBody>
      <Heading title="비밀번호 초기화" />
      <ResetPasswordClient />
    </BlackSideBody>
  );
};

export default ResetPassword;
