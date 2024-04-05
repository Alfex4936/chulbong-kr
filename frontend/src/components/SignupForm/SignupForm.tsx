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

import { useTranslation } from 'react-i18next';
import '../../i18n';

const SignupForm = () => {
  const { t } = useTranslation();

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
  const [viewTimer, setViewTimer] = useState(true);
  const [successMessage, setSuccessMessage] = useState("");

  useEffect(() => {
    toastState.close();
    toastState.setToastText("");
  }, []);

  const handleSubmit = () => {
    toastState.close();
    let isValid = true;

    if (!validattionComplete) {
      setCodeError(t("signup.completeVerification"));
      isValid = false;
    } else {
      setCodeError("");
    }

    if (nameInput.value === "") {
      setNameError(t("signup.enterNickname"));
      isValid = false;
    } else {
      setNameError("");
    }

    if (emailInput.value === "") {
      setEmailError(t("login.enterEmail"));
      isValid = false;
    } else if (!emailValidate(emailInput.value)) {
      setEmailError(t("login.emailInvalidFormat"));
      isValid = false;
    } else {
      setEmailError("");
    }

    if (passwordInput.value === "") {
      setPasswordError(t("login.enterPassword"));
      isValid = false;
    } else if (!passwordValidate(passwordInput.value)) {
      setPasswordError(t("login.passwordRequirements"));
      isValid = false;
    } else {
      setPasswordError("");
    }

    if (verifyPasswordInput.value === "") {
      setVerifyPasswordError(t("login.enterPassword"));
      isValid = false;
    } else if (passwordInput.value !== verifyPasswordInput.value) {
      setVerifyPasswordError(t("signup.confirmPassword"));
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
          toastState.setToastText(t("signup.signupComplete"));
          setSuccessMessage(t("signup.signupComplete") + "!");
          toastState.open();
          modalState.close();
          modalState.openLogin();
        })
        .catch((error) => {
          setLoading(false);
          if (error.response.status === 409) {
            setSigninError(t("signup.alreadyRegistered"));
          } else {
            setSigninError(t("signup.tryAgainLater"));
          }
          console.log(error);
        });
    }
  };

  const handleGetCode = async () => {
    setStartTimer(false);
    if (emailInput.value === "") {
      setEmailError(t("login.enterEmail"));
      return;
    } else if (!emailValidate(emailInput.value)) {
      setEmailError(t("login.emailInvalidFormat"));
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
        if (error.response?.status === 409) {
          setEmailError(t("login.emailAlreadyRegistered"));
        } else {
          setEmailError(t("login.error"));
        }
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
      setViewTimer(false);
      console.log(result);
    } catch (error) {
      if (isAxiosError(error)) {
        console.log(error.response?.status);
        if (error.response?.status === 400) {
          setCodeError(t("login.invalidCode"));
        }
      }
    } finally {
      setValidateCodeLoading(false);
    }
  };

  return (
    <form>
      <Styled.FormTitle>{t("signup.signup")}</Styled.FormTitle>
      <Styled.InputWrap>
        <Input
          type="text"
          id="name"
          data-testid="name"
          placeholder={t("signup.nickname")}
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
          data-testid="email"
          placeholder={t("map.auth.email")}
          buttonText={
            getCodeLoading ? (
              <CircularProgress size={15} sx={{ color: "#fff" }} />
            ) : getCodeComplete ? (
              t("signup.requestAgain")
            ) : (
              t("signup.requestVerification")
            )
          }
          value={emailInput.value}
          onChange={(e) => {
            emailInput.onChange(e);
            setEmailError("");
          }}
          onClickFn={handleGetCode}
        />
        <Styled.ErrorBox data-testid="email-error">
          {emailError}
        </Styled.ErrorBox>
      </Styled.InputWrap>
      {getCodeComplete && (
        <Styled.InputWrap>
          <Input
            maxLength={6}
            theme="button"
            type="number"
            id="code"
            data-testid="code"
            placeholder={t("signup.verificationCode")}
            buttonText={
              validateCodeLoading ? (
                <CircularProgress size={15} sx={{ color: "#fff" }} />
              ) : validattionComplete ? (
                <CheckIcon />
              ) : (
                t("signup.verify")
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
            {viewTimer && (
              <Styled.TimerContainer>
                <CertificationCount
                  start={startTimer}
                  setStart={setStartTimer}
                />
              </Styled.TimerContainer>
            )}
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
          data-testid="password"
          placeholder={t("map.auth.password")}
          value={passwordInput.value}
          onChange={(e) => {
            passwordInput.onChange(e);
            setPasswordError("");
          }}
        />
        <Styled.ErrorBox data-testid="password-error">
          {passwordError}
        </Styled.ErrorBox>
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
          data-testid="verify-password"
          placeholder={t("signup.confirmPasswordLabel")}
          value={verifyPasswordInput.value}
          onChange={(e) => {
            verifyPasswordInput.onChange(e);
            setVerifyPasswordError("");
          }}
        />
        <Styled.ErrorBox data-testid="verify-password-error">
          {verifyPasswordError}
        </Styled.ErrorBox>
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
        data-testid="signup-button"
        disabled={loading}
      >
        {loading ? (
          <CircularProgress size={20} sx={{ color: "#fff" }} />
        ) : (
          t("signup.signup")
        )}
      </Button>
      <div data-testid="signup-success">{successMessage}</div>
    </form>
  );
};

export default SignupForm;
