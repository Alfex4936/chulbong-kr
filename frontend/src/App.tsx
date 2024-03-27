import VisibilityIcon from "@mui/icons-material/Visibility";
import VisibilityOffIcon from "@mui/icons-material/VisibilityOff";
import CircularProgress from "@mui/material/CircularProgress";
import { AxiosError } from "axios";
import { Suspense, lazy, useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { Flip, ToastContainer, toast } from "react-toastify";
import * as Styled from "./App.style";
import ActionButton from "./components/ActionButton/ActionButton";
import Input from "./components/Input/Input";
import LoginFormSkeleton from "./components/LoginForm/LoginFormSkeleton";
import Map from "./components/Map/Map";
import BasicModal from "./components/Modal/Modal";
import SignupFormSkeleton from "./components/SignupForm/SignupFormSkeleton";
import useLogout from "./hooks/mutation/auth/useLogout";
import useResetPassword from "./hooks/mutation/auth/useResetPassword";
import useGetMyInfo from "./hooks/query/user/useGetMyInfo";
import useInput from "./hooks/useInput";
import useModalStore from "./store/useModalStore";
import useOnBoardingStore from "./store/useOnBoardingStore";
import useToastStore from "./store/useToastStore";
import useUserStore from "./store/useUserStore";
import { nanoid } from "nanoid";
import useChatIdStore from "./store/useChatIdStore";

import "react-toastify/dist/ReactToastify.css";

const LoginForm = lazy(() => import("./components/LoginForm/LoginForm"));
const SignupForm = lazy(() => import("./components/SignupForm/SignupForm"));

const App = () => {
  const modalState = useModalStore();
  const toastState = useToastStore();
  const userState = useUserStore();
  const onBoardingState = useOnBoardingStore();
  const cidState = useChatIdStore();

  const location = useLocation();
  const navigate = useNavigate();

  const query = new URLSearchParams(location.search);
  const token = query.get("token");
  const email = query.get("email");

  const passwordInput = useInput("");
  const {
    mutate,
    isError: resetError,
    isPending,
    isSuccess,
  } = useResetPassword(token as string, passwordInput.value);
  const { mutate: logout } = useLogout();

  const { isError, error } = useGetMyInfo();

  const [viewPassword, setViewPassword] = useState(false);

  const [changePasswordModal, setChangePasswordModal] = useState(false);

  useEffect(() => {
    const lastVisit = localStorage.getItem("lastVisit");
    const now = new Date();

    if (lastVisit) {
      const lastVisitDate = new Date(lastVisit);
      const daysDifference = Math.floor(
        (now.getTime() - lastVisitDate.getTime()) / (1000 * 60 * 60 * 24)
      );

      if (daysDifference >= 4) {
        if (!userState.ka.h) onBoardingState.open();
      }
    } else {
      if (!userState.ka.h) onBoardingState.open();
    }

    localStorage.setItem("lastVisit", now.toISOString());

    const setId = () => {
      if (!cidState.cid) {
        cidState.setId(nanoid());
      }
    };

    setId();

    window.addEventListener("storage", setId);

    return () => {
      window.removeEventListener("storage", setId);
    };
  }, []);

  useEffect(() => {
    if (token && email) {
      setChangePasswordModal(true);
    }
  }, [token, email]);

  useEffect(() => {
    const handleLogout = () => {
      logout();
      userState.resetUser();
    };

    if (isError) {
      if (error instanceof AxiosError) {
        handleLogout();
      }
    }
  }, [isError, error]);

  const notify = () => toast(toastState.toastText);

  useEffect(() => {
    if (toastState.isToast) {
      notify();
    }
  }, [toastState.isToast]);

  useEffect(() => {
    if (isSuccess) {
      setChangePasswordModal(false);
      toastState.setToastText("비밀번호 변경 완료");
      toastState.open();
      navigate("/");
    }
  }, [isSuccess]);

  const handleChangePassword = () => {
    mutate();
  };

  return (
    <div>
      <Map />
      {modalState.loginModal && (
        <BasicModal>
          <Suspense fallback={<LoginFormSkeleton />}>
            <LoginForm />
          </Suspense>
        </BasicModal>
      )}
      {modalState.signupModal && (
        <BasicModal>
          <Suspense fallback={<SignupFormSkeleton />}>
            <SignupForm />
          </Suspense>
        </BasicModal>
      )}

      {changePasswordModal && (
        <BasicModal setState={setChangePasswordModal}>
          <p
            style={{
              margin: "1rem",
              fontSize: "1.2rem",
            }}
          >
            비밀번호 변경
          </p>
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
            }}
          />
          {resetError && (
            <Styled.ErrorBox>잠시 후 다시 시도해주세요!</Styled.ErrorBox>
          )}

          <Styled.ChangePasswordButtonsWrap>
            <ActionButton bg="black" onClick={handleChangePassword}>
              {isPending ? (
                <CircularProgress size={20} sx={{ color: "#fff" }} />
              ) : (
                "변경"
              )}
            </ActionButton>
            <ActionButton
              bg="gray"
              onClick={() => {
                setChangePasswordModal(false);
              }}
            >
              취소
            </ActionButton>
          </Styled.ChangePasswordButtonsWrap>
        </BasicModal>
      )}

      <ToastContainer
        position="top-right"
        autoClose={1000}
        hideProgressBar={false}
        newestOnTop={false}
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover
        theme="dark"
        transition={Flip}
        style={{ width: "96%", maxWidth: "300px" }}
      />
    </div>
  );
};

export default App;
