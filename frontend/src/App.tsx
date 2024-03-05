import { Suspense, lazy, useEffect } from "react";
import { Flip, ToastContainer, toast } from "react-toastify";
import LoginFormSkeleton from "./components/LoginForm/LoginFormSkeleton";
import Map from "./components/Map/Map";
import BasicModal from "./components/Modal/Modal";
import SignupFormSkeleton from "./components/SignupForm/SignupFormSkeleton";
import useModalStore from "./store/useModalStore";
import useToastStore from "./store/useToastStore";

import "react-toastify/dist/ReactToastify.css";
import useUserStore from "./store/useUserStore";
import useGetMyInfo from "./hooks/query/user/useGetMyInfo";
import logout from "./api/auth/logout";
import { AxiosError } from "axios";

const LoginForm = lazy(() => import("./components/LoginForm/LoginForm"));
const SignupForm = lazy(() => import("./components/SignupForm/SignupForm"));

const App = () => {
  const modalState = useModalStore();
  const toastState = useToastStore();
  const userState = useUserStore();

  const { isError, error } = useGetMyInfo();

  useEffect(() => {
    const handleLogout = async () => {
      try {
        const result = await logout();
        userState.resetUser();
        console.log(result);
      } catch (error) {
        userState.resetUser();
        console.log(error);
      }
    };

    if (isError) {
      if (error instanceof AxiosError) {
        handleLogout();
        console.log(error.response?.status);
      } else {
        console.error(error);
      }
    }
  }, [isError]);

  const notify = () => toast(toastState.toastText);

  useEffect(() => {
    if (toastState.isToast) {
      notify();
    }
  }, [toastState.isToast]);

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
      />
    </div>
  );
};

export default App;
