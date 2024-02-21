import { useEffect } from "react";
import { Flip, ToastContainer, toast } from "react-toastify";
import LoginForm from "./components/LoginForm/LoginForm";
import Map from "./components/Map/Map";
import BasicModal from "./components/Modal/Modal";
import SignupForm from "./components/SignupForm/SignupForm";
import useModalStore from "./store/useModalStore";
import useToastStore from "./store/useToastStore";

import "react-toastify/dist/ReactToastify.css";

const App = () => {
  const modalState = useModalStore();
  const toastState = useToastStore();

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
          <LoginForm />
        </BasicModal>
      )}
      {modalState.signupModal && (
        <BasicModal>
          <SignupForm />
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
