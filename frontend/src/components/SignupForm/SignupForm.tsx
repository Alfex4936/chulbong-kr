import * as Styled from "./SignupForm.style";
import Input from "../Input/Input";
import { Button } from "@mui/material";

const SignupForm = () => {
  const handleSubmit = () => {
    console.log("회원가입");
  };

  return (
    <Styled.HiddenBox>
      <Styled.FormWrap>
        <Styled.FormTitle>회원가입</Styled.FormTitle>
        <Styled.InputWrap>
          <label htmlFor="name">닉네임</label>
          <Input type="text" id="name" />
        </Styled.InputWrap>
        <Styled.InputWrap>
          <label htmlFor="email">이메일</label>
          <Input type="email" id="email" />
        </Styled.InputWrap>
        <Styled.InputWrap>
          <label htmlFor="password">비밀번호</label>
          <Input type="password" id="password" />
        </Styled.InputWrap>
        <Styled.InputWrap>
          <label htmlFor="verify-password">비밀번호 확인</label>
          <Input type="password" id="verify-password" />
        </Styled.InputWrap>
        <Button
          onClick={handleSubmit}
          sx={{
            color: "#fff",
            width: "100%",
            backgroundColor: "#333",
            margin: "1.5rem 0",
            "&:hover": {
              backgroundColor: "#555",
            },
          }}
        >
          회원가입
        </Button>
      </Styled.FormWrap>
    </Styled.HiddenBox>
  );
};

export default SignupForm;
