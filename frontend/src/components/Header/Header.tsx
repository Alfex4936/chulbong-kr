import * as Styled from "./Header.style";
import useModalStore from "../../store/useModalStore";
import { Button } from "@mui/material";

const Header = () => {
  const modalState = useModalStore();

  const handleOpen = () => {
    modalState.openLogin();
    console.log("로그인 클릭");
  };

  return (
    <Styled.HeaderContainer>
      <div>철봉</div>
      {/* <button onClick={handleOpen}>로그인/회원가입</button> */}
      <Button onClick={handleOpen} sx={{ color: "#333" }}>
        로그인/회원가입
      </Button>
    </Styled.HeaderContainer>
  );
};

export default Header;
