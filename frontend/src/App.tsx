import Map from "./components/Map/Map";
import Header from "./components/Header/Header";
import LoginForm from "./components/LoginForm/LoginForm";
import BasicModal from "./components/Modal/Modal";
import useModalStore from "./store/useModalStore";
import SignupForm from "./components/SignupForm/SignupForm";

const App = () => {
  const modalState = useModalStore();

  return (
    <div>
      <Header />
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
    </div>
  );
};

export default App;
