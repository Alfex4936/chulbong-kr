import { Button, CircularProgress } from "@mui/material";
import { useEffect, useState } from "react";
import signin from "../../api/auth/signin";
import useInput from "../../hooks/useInput";
import useModalStore from "../../store/useModalStore";
import useToastStore from "../../store/useToastStore";
import emailValidate from "../../utils/emailValidate";
import passwordValidate from "../../utils/passwordValidate";
import Input from "../Input/Input";
import * as Styled from "./SignupForm.style";

const SignupForm = () => {
  const modalState = useModalStore();
  const toastState = useToastStore();

  const nameInput = useInput("");
  const emailInput = useInput("");
  const passwordInput = useInput("");
  const verifyPasswordInput = useInput("");

  const [nameError, setNameError] = useState("");
  const [emailError, setEmailError] = useState("");
  const [passwordError, setPasswordError] = useState("");
  const [verifyPasswordError, setVerifyPasswordError] = useState("");
  const [signinError, setSigninError] = useState("");

  const [loading, setLoading] = useState(false);

  useEffect(() => {
    toastState.close();
    toastState.setToastText("");
  }, []);

  const handleSubmit = () => {
    toastState.close();
    let isValid = true;

    if (nameInput.value === "") {
      setNameError("닉네임을 입력해 주세요");
      isValid = false;
    } else {
      setNameError("");
    }

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
      setPasswordError("비밀번호를 입력해 주세요");
      isValid = false;
    } else if (!passwordValidate(passwordInput.value)) {
      setPasswordError("특수문자 포함 8 ~ 20자 사이로 입력해 주세요.");
      isValid = false;
    } else {
      setPasswordError("");
    }

    if (verifyPasswordInput.value === "") {
      setVerifyPasswordError("비밀번호를 입력해 주세요.");
      isValid = false;
    } else if (passwordInput.value !== verifyPasswordInput.value) {
      setVerifyPasswordError("비밀번호를 확인해 주세요.");
      isValid = false;
    } else {
      setVerifyPasswordError("");
    }

    if (isValid) {
      setLoading(true);
      signin({
        username: nameInput.value,
        email: emailInput.value,
        password: passwordInput.value,
      })
        .then(() => {
          toastState.setToastText("회원 가입 완료");
          toastState.open();
          modalState.close();
          modalState.openLogin();
        })
        .catch((error) => {
          setLoading(false);
          if (error.response.status === 409) {
            setSigninError("이미 등록된 회원입니다.");
          } else {
            setSigninError("잠시 후 다시 시도해 주세요.");
          }
          console.log(error);
        });
    }
  };

  return (
    <form>
      <Styled.FormTitle>회원가입</Styled.FormTitle>
      <Styled.InputWrap>
        <Input
          type="text"
          id="name"
          placeholder="닉네임"
          value={nameInput.value}
          onChange={(e) => {
            nameInput.onChange(e);
            setNameError("");
          }}
        />
        <Styled.ErrorBox>{nameError}</Styled.ErrorBox>
      </Styled.InputWrap>
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
      <Styled.InputWrap>
        <Input
          type="password"
          id="verify-password"
          placeholder="비밀번호 확인"
          value={verifyPasswordInput.value}
          onChange={(e) => {
            verifyPasswordInput.onChange(e);
            setVerifyPasswordError("");
          }}
        />
        <Styled.ErrorBox>{verifyPasswordError}</Styled.ErrorBox>
      </Styled.InputWrap>
      <Styled.ErrorBox>{signinError}</Styled.ErrorBox>
      <Button
        onClick={handleSubmit}
        sx={{
          color: "#fff",
          width: "100%",
          height: "40px",
          backgroundColor: "#333",
          margin: "1rem 0",
          "&:hover": {
            backgroundColor: "#555",
          },
        }}
        disabled={loading}
      >
        {loading ? (
          <CircularProgress size={20} sx={{ color: "#fff" }} />
        ) : (
          "회원가입"
        )}
      </Button>
    </form>
  );
};

export default SignupForm;
