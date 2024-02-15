import { Button } from "@mui/material";
import useModalStore from "../../store/useModalStore";
import * as Styled from "./Header.style";
import useUserStore from "../../store/useUserStore";

const Header = () => {
  const modalState = useModalStore();
  const userState = useUserStore();

  const handleOpen = () => {
    modalState.openLogin();
    console.log("로그인 클릭");
  };
  const handleLogin = () => {
    userState.resetUser();
    console.log("로그아웃");
  };

  return (
    <Styled.HeaderContainer>
      <div>철봉</div>
      {userState.user.token === "" ? (
        <Button onClick={handleOpen} sx={{ color: "#333" }}>
          로그인/회원가입
        </Button>
      ) : (
        <Button onClick={handleLogin} sx={{ color: "#333" }}>
          {userState.user.user.username}
        </Button>
      )}
    </Styled.HeaderContainer>
  );
};

export default Header;
