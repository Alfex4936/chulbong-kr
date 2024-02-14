import * as Styled from "./LoginForm.style";
import Input from "../Input/Input";

const LoginForm = () => {
  const handleSubmit = () => {
    console.log("로그인");
  };

  return (
    <div>
      <Styled.FormTitle>로그인</Styled.FormTitle>
      <Styled.InputWrap>
        <label htmlFor="email">이메일</label>
        <Input type="email" id="email" />
      </Styled.InputWrap>
      <Styled.InputWrap>
        <label htmlFor="password">비밀번호</label>
        <Input type="password" id="password" />
      </Styled.InputWrap>
      <button onClick={handleSubmit}>로그인</button>
    </div>
  );
};

export default LoginForm;
