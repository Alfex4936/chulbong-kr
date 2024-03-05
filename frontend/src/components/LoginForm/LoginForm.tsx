import VisibilityIcon from "@mui/icons-material/Visibility";
import VisibilityOffIcon from "@mui/icons-material/VisibilityOff";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import { useState } from "react";
import login from "../../api/auth/login";
import useInput from "../../hooks/useInput";
import useModalStore from "../../store/useModalStore";
import useUserStore from "../../store/useUserStore";
import emailValidate from "../../utils/emailValidate";
import passwordValidate from "../../utils/passwordValidate";
import Input from "../Input/Input";
import * as Styled from "./LoginForm.style";

const LoginForm = () => {
  const modalState = useModalStore();
  const userState = useUserStore();

  const emailInput = useInput("");
  const passwordInput = useInput("");

  const [emailError, setEmailError] = useState("");
  const [passwordError, setPasswordError] = useState("");
  const [loginError, setLoginError] = useState("");

  const [viewPassword, setViewPassword] = useState(false);

  const [loading, setLoading] = useState(false);

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
      setLoading(true);
      login({
        email: emailInput.value,
        password: passwordInput.value,
      })
        .then((res) => {
          setLoginError("");
          userState.setUser(res);
          modalState.close();
        })
        .catch((error) => {
          console.log(error);
          setLoginError("유요하지 않은 회원 정보입니다.");
        })
        .finally(() => {
          setLoading(false);
        });
    }
  };

  const handleClickEmailSignin = () => {
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
          theme="icon"
          icon={
            viewPassword ? (
              <VisibilityIcon fontSize="small" />
            ) : (
              <VisibilityOffIcon fontSize="small" />
            )
          }
          onClickFn={() => {
            setViewPassword((prev) => !prev);
          }}
          type={viewPassword ? "text" : "password"}
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
      <Styled.ErrorBox>{loginError}</Styled.ErrorBox>
      <Button
        onClick={handleSubmit}
        sx={{
          color: "#fff",
          width: "100%",
          height: "40px",
          backgroundColor: "#333",
          marginTop: "1rem",
          "&:hover": {
            backgroundColor: "#555",
          },
        }}
        disabled={loading}
      >
        {loading ? (
          <CircularProgress size={20} sx={{ color: "#fff" }} />
        ) : (
          "로그인"
        )}
      </Button>
      <div style={{ marginTop: "1rem" }}>
        <Styled.SignupButtonWrap>
          <p>계정이 없으신가요?</p>
          <Styled.SigninLinkButton onClick={handleClickEmailSignin}>
            이메일로 회원가입 하기
          </Styled.SigninLinkButton>
        </Styled.SignupButtonWrap>
        <Styled.SignupButtonWrap>
          <p>비밀번호를 잊어버리셨나요?</p>
          <Styled.SigninLinkButton
            onClick={() => {
              modalState.close();
              modalState.openPassword();
            }}
          >
            비밀번호 변경하기
          </Styled.SigninLinkButton>
        </Styled.SignupButtonWrap>
      </div>
    </form>
  );
};

export default LoginForm;
