import { Button } from "@mui/material";
import { useState } from "react";
import useInput from "../../hooks/useInput";
import useModalStore from "../../store/useModalStore";
import emailValidate from "../../utils/emailValidate";
import passwordValidate from "../../utils/passwordValidate";
import Input from "../Input/Input";
import * as Styled from "./LoginForm.style";

const LoginForm = () => {
  const modalState = useModalStore();

  const emailInput = useInput("");
  const passwordInput = useInput("");

  const [emailError, setEmailError] = useState("");
  const [passwordError, setPasswordError] = useState("");

  const handleSubmit = () => {
    let isValid = true;

    if (emailInput.value === "") {
      setEmailError("이메일을 입력해 주세요");
      isValid = false;
    } else if (!emailValidate(emailInput.value)) {
      setEmailError("이메일 형식이 아닙니다.");
      isValid = false;
    } else {
      setEmailError("");
    }

    if (passwordInput.value === "") {
      setPasswordError("비밀번호 입력해 주세요");
      isValid = false;
    } else if (!passwordValidate(passwordInput.value)) {
      setPasswordError("특수문자 포함 8 ~ 20자 사이로 입력해 주세요.");
      isValid = false;
    } else {
      setPasswordError("");
    }

    if (isValid) {
      console.log({
        email: emailInput.value,
        password: passwordInput.value,
      });
    }
  };

  const handleClickEmailSignin = () => {
    console.log("이메일 회원가입 하기");
    modalState.openSignup();
    modalState.closeLogin();
  };

  return (
    <form>
      <Styled.FormTitle>로그인</Styled.FormTitle>
      <Styled.InputWrap>
        <Input
          type="email"
          id="email"
          placeholder="이메일"
          value={emailInput.value}
          onChange={(e) => {
            emailInput.onChange(e);
            setEmailError("");
          }}
        />
        <Styled.ErrorBox>{emailError}</Styled.ErrorBox>
      </Styled.InputWrap>
      <Styled.InputWrap>
        <Input
          type="password"
          id="password"
          placeholder="비밀번호"
          value={passwordInput.value}
          onChange={(e) => {
            passwordInput.onChange(e);
            setPasswordError("");
          }}
        />
        <Styled.ErrorBox>{passwordError}</Styled.ErrorBox>
      </Styled.InputWrap>
      <Button
        onClick={handleSubmit}
        sx={{
          color: "#fff",
          width: "100%",
          backgroundColor: "#333",
          marginTop: "1rem",
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
    </form>
  );
};

export default LoginForm;
