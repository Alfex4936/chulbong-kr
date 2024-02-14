import * as Styled from "./LoginForm.style";
import Input from "../Input/Input";
import { Button } from "@mui/material";
import useModalStore from "../../store/useModalStore";

const LoginForm = () => {
  const modalState = useModalStore();

  const handleSubmit = () => {
    console.log("로그인");
  };

  const handleClickEmailSignin = () => {
    console.log("이메일 회원가입 하기");
    modalState.openSignup();
    modalState.closeLogin();
  };

  return (
    <Styled.HiddenBox>
      <Styled.FormWrap>
        <Styled.FormTitle>로그인</Styled.FormTitle>
        <Styled.InputWrap>
          <label htmlFor="email">이메일</label>
          <Input type="email" id="email" />
        </Styled.InputWrap>
        <Styled.InputWrap>
          <label htmlFor="password">비밀번호</label>
          <Input type="password" id="password" />
        </Styled.InputWrap>
        <Button
          onClick={handleSubmit}
          sx={{
            color: "#fff",
            width: "100%",
            backgroundColor: "#333",
            marginTop: "1.5rem",
            "&:hover": {
              backgroundColor: "#555",
            },
          }}
        >
          로그인
        </Button>
        <Styled.SignupButtonWrap>
          <p>계정이 없으신가요?</p>
          <Styled.SigninLinkButton onClick={handleClickEmailSignin}>
            이메일로 회원가입 하기
          </Styled.SigninLinkButton>
        </Styled.SignupButtonWrap>
      </Styled.FormWrap>
    </Styled.HiddenBox>
  );
};

export default LoginForm;
