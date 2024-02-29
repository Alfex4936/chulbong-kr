import * as Styled from "./LoginFormSkeleton.style";

const LoginFormSkeleton = () => {
  return (
    <div>
      <Styled.TitleSkeleton />
      <Styled.InputSkeleton />
      <Styled.InputSkeleton />
      <Styled.ButtonSkeleton />
      <Styled.SigninSkeleton />
    </div>
  );
};

export default LoginFormSkeleton;
