import { Button } from "@mui/material";
import useModalStore from "../../store/useModalStore";
import * as Styled from "./Header.style";
import useUserStore from "../../store/useUserStore";

const Header = () => {
  const modalState = useModalStore();
  const userState = useUserStore();

  const handleOpen = () => {
    modalState.openLogin();
  };
  const handleLogout = () => {
    userState.resetUser();
  };

  return (
    <Styled.HeaderContainer>
      <div>철봉</div>
      {userState.user.token === "" ? (
        <Button onClick={handleOpen} sx={{ color: "#333" }}>
          로그인/회원가입
        </Button>
      ) : (
        <Button onClick={handleLogout} sx={{ color: "#333" }}>
          {userState.user.user.username}
        </Button>
      )}
    </Styled.HeaderContainer>
  );
};

export default Header;
