import CheckIcon from "@mui/icons-material/Check";
import VisibilityIcon from "@mui/icons-material/Visibility";
import VisibilityOffIcon from "@mui/icons-material/VisibilityOff";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import { isAxiosError } from "axios";
import { useEffect, useState } from "react";
import sendVerifyCode from "../../api/auth/sendVerifyCode";
import signin from "../../api/auth/signin";
import verifyCode from "../../api/auth/verifyCode";
import useInput from "../../hooks/useInput";
import useModalStore from "../../store/useModalStore";
import useToastStore from "../../store/useToastStore";
import emailValidate from "../../utils/emailValidate";
import passwordValidate from "../../utils/passwordValidate";
import Input from "../Input/Input";
import CertificationCount from "./CertificationCount";
import * as Styled from "./SignupForm.style";

const SignupForm = () => {
  const modalState = useModalStore();
  const toastState = useToastStore();

  const nameInput = useInput("");
  const emailInput = useInput("");
  const codeInput = useInput("");
  const passwordInput = useInput("");
  const verifyPasswordInput = useInput("");

  const [nameError, setNameError] = useState("");
  const [emailError, setEmailError] = useState("");
  const [codeError, setCodeError] = useState("");
  const [passwordError, setPasswordError] = useState("");
  const [verifyPasswordError, setVerifyPasswordError] = useState("");
  const [signinError, setSigninError] = useState("");

  const [loading, setLoading] = useState(false);

  const [getCodeLoading, setGetCodeLoading] = useState(false);
  const [getCodeComplete, setGetCodeComplete] = useState(false);

  const [validateCodeLoading, setValidateCodeLoading] = useState(false);
  const [validattionComplete, setValidattionComplete] = useState(false);

  const [viewPassword, setViewPassword] = useState(false);
  const [viewVerifyPassword, setViewVerifyPassword] = useState(false);

  const [startTimer, setStartTimer] = useState(false);

  useEffect(() => {
    toastState.close();
    toastState.setToastText("");
  }, []);

  const handleSubmit = () => {
    toastState.close();
    let isValid = true;

    if (!validattionComplete) {
      setCodeError("인증을 완료해 주세요");
      isValid = false;
    } else {
      setCodeError("");
    }

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

  const handleGetCode = async () => {
    setStartTimer(false);
    if (emailInput.value === "") {
      setEmailError("이메일을 입력해 주세요");
      return;
    } else if (!emailValidate(emailInput.value)) {
      setEmailError("이메일 형식이 아닙니다.");
      return;
    }

    setGetCodeLoading(true);

    try {
      const result = await sendVerifyCode(emailInput.value);
      setStartTimer(true);
      setGetCodeComplete(true);
      setEmailError("");
      console.log(result);
    } catch (error) {
      if (isAxiosError(error)) {
        console.log(error);
        setEmailError("에러");
      }
    } finally {
      setGetCodeLoading(false);
    }
  };

  const handleSubmitCode = async () => {
    setValidateCodeLoading(true);
    try {
      const result = await verifyCode({
        email: emailInput.value,
        code: codeInput.value,
      });

      setCodeError("");
      setValidattionComplete(true);
      console.log(result);
    } catch (error) {
      if (isAxiosError(error)) {
        console.log(error.response?.status);
        if (error.response?.status === 400) {
          setCodeError("유효하지 않은 코드입니다. 다시 시도해 주세요.");
        }
      }
    } finally {
      setValidateCodeLoading(false);
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
          theme="button"
          type="email"
          id="email"
          placeholder="이메일"
          buttonText={
            getCodeLoading ? (
              <CircularProgress size={15} sx={{ color: "#fff" }} />
            ) : getCodeComplete ? (
              "다시 요청"
            ) : (
              "인증 요청"
            )
          }
          value={emailInput.value}
          onChange={(e) => {
            emailInput.onChange(e);
            setEmailError("");
          }}
          onClickFn={handleGetCode}
        />
        <Styled.ErrorBox>{emailError}</Styled.ErrorBox>
      </Styled.InputWrap>
      {getCodeComplete && (
        <Styled.InputWrap>
          <Input
            maxLength={6}
            theme="button"
            type="number"
            id="code"
            placeholder="인증번호"
            buttonText={
              validateCodeLoading ? (
                <CircularProgress size={15} sx={{ color: "#fff" }} />
              ) : validattionComplete ? (
                <CheckIcon />
              ) : (
                "인증 확인"
              )
            }
            value={codeInput.value}
            onChange={(e) => {
              if (/^\d{0,6}$/.test(e.target.value)) {
                codeInput.onChange(e);
              }
            }}
            onClickFn={handleSubmitCode}
          />
          <div style={{ display: "flex" }}>
            <Styled.ErrorBox>{codeError}</Styled.ErrorBox>
            <div style={{ flexGrow: "1" }} />
            <Styled.TimerContainer>
              <CertificationCount start={startTimer} setStart={setStartTimer} />
            </Styled.TimerContainer>
          </div>
        </Styled.InputWrap>
      )}

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
      <Styled.InputWrap>
        <Input
          theme="icon"
          icon={
            viewVerifyPassword ? (
              <VisibilityIcon fontSize="small" />
            ) : (
              <VisibilityOffIcon fontSize="small" />
            )
          }
          onClickFn={() => {
            setViewVerifyPassword((prev) => !prev);
          }}
          type={viewVerifyPassword ? "text" : "password"}
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
