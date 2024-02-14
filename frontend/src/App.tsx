import Map from "./components/Map/Map";
import Header from "./components/Header/Header";
import LoginForm from "./components/LoginForm/LoginForm";
import BasicModal from "./components/Modal/Modal";
import useModalStore from "./store/useModalStore";

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
    </div>
  );
};

export default App;
