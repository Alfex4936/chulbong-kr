import * as Styled from "./Header.style";
import useModalStore from "../../store/useModalStore";

const Header = () => {
  const modalState = useModalStore();

  const handleOpen = () => {
    modalState.openLogin();
    console.log("로그인 클릭");
  };

  return (
    <Styled.HeaderContainer>
      <div>철봉</div>
      <button onClick={handleOpen}>로그인/회원가입</button>
    </Styled.HeaderContainer>
  );
};

export default Header;
