import VisibilityIcon from "@mui/icons-material/Visibility";
import VisibilityOffIcon from "@mui/icons-material/VisibilityOff";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import { useEffect, useState } from "react";
import useSignin from "../../hooks/mutation/auth/useSignin";
import useInput from "../../hooks/useInput";
import useModalStore from "../../store/useModalStore";
import useUserStore from "../../store/useUserStore";
import emailValidate from "../../utils/emailValidate";
import passwordValidate from "../../utils/passwordValidate";
import Input from "../Input/Input";
import * as Styled from "./LoginForm.style";

import { useTranslation } from 'react-i18next';
import '../../i18n';

const LoginForm = () => {
  const { t } = useTranslation();

  const modalState = useModalStore();
  const userState = useUserStore();

  const emailInput = useInput("");
  const passwordInput = useInput("");

  const { mutateAsync: login } = useSignin();

  const [emailError, setEmailError] = useState("");
  const [passwordError, setPasswordError] = useState("");
  const [loginError, setLoginError] = useState("");

  const [viewPassword, setViewPassword] = useState(false);

  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const handleSubmit = async () => {
      let isValid = true;

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

      if (isValid) {
        setLoading(true);
        try {
          await login({
            email: emailInput.value,
            password: passwordInput.value,
          });
          setLoginError("");
          modalState.close();
          userState.setLogin();
        } catch (error) {
          setLoginError(t("login.invalidCredentials"));
          userState.resetUser();
        } finally {
          setLoading(false);
        }
      }
    };

    const handleKeyDownClose = (event: KeyboardEvent) => {
      if (event.key === "Enter") {
        if (loading) return;
        handleSubmit();
      }
    };

    window.addEventListener("keydown", handleKeyDownClose);

    return () => {
      window.removeEventListener("keydown", handleKeyDownClose);
    };
  }, [emailInput.value, passwordInput.value, loading]);

  const handleSubmit = async () => {
    let isValid = true;

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

    if (isValid) {
      setLoading(true);
      try {
        await login({
          email: emailInput.value,
          password: passwordInput.value,
        });
        setLoginError("");
        modalState.close();
        userState.setLogin();
      } catch (error) {
        console.log(error);
        setLoginError(t("login.invalidCredentials"));
        userState.resetUser();
      } finally {
        setLoading(false);
      }
    }
  };

  const handleClickEmailSignin = () => {
    modalState.openSignup();
    modalState.closeLogin();
  };

  return (
    <form>
      <Styled.FormTitle>{t("map.auth.login")}</Styled.FormTitle>
      <Styled.InputWrap>
        <Input
          type="email"
          id="email"
          placeholder={t("map.auth.email")}
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
          placeholder={t("map.auth.password")}
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
          t("map.auth.login")
        )}
      </Button>
      <div style={{ marginTop: "1rem" }}>
        <Styled.SignupButtonWrap>
          <p>{t("login.noAccount")}</p>
          <Styled.SigninLinkButton onClick={handleClickEmailSignin}>
            {t("login.signUpWithEmail")}
          </Styled.SigninLinkButton>
        </Styled.SignupButtonWrap>
        <Styled.SignupButtonWrap>
          <p>{t("login.forgotPassword")}</p>
          <Styled.SigninLinkButton
            onClick={() => {
              modalState.close();
              modalState.openPassword();
            }}
          >
            {t("login.changePassword")}
          </Styled.SigninLinkButton>
        </Styled.SignupButtonWrap>
      </div>
    </form>
  );
};

export default LoginForm;
