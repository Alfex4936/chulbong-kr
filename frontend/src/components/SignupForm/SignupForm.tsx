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
          <Input type="text" id="name" placeholder="닉네임" />
        </Styled.InputWrap>
        <Styled.InputWrap>
          <Input type="email" id="email" placeholder="이메일" />
        </Styled.InputWrap>
        <Styled.InputWrap>
          <Input type="password" id="password" placeholder="비밀번호" />
        </Styled.InputWrap>
        <Styled.InputWrap>
          <Input
            type="password"
            id="verify-password"
            placeholder="비밀번호 확인"
          />
        </Styled.InputWrap>
        <Button
          onClick={handleSubmit}
          sx={{
            color: "#fff",
            width: "100%",
            backgroundColor: "#333",
            margin: "1rem 0",
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
