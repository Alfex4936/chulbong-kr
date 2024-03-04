import CloseIcon from "@mui/icons-material/Close";
import EditIcon from "@mui/icons-material/Edit";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { useState } from "react";
import logout from "../../api/auth/logout";
import useUserStore from "../../store/useUserStore";
import ActionButton from "../ActionButton/ActionButton";
import * as Styled from "./MyInfoModal.style";

interface Props {
  setMyInfoModal: React.Dispatch<React.SetStateAction<boolean>>;
}

const MyInfoModal = ({ setMyInfoModal }: Props) => {
  const userState = useUserStore();

  const [logoutLoading, setLogoutLoading] = useState(false);

  const handleLogout = async () => {
    setLogoutLoading(true);
    try {
      const result = await logout();
      userState.resetUser();
      setMyInfoModal(false);
      console.log(result);
    } catch (error) {
      userState.resetUser();
      setMyInfoModal(false);
      console.log(error);
    } finally {
      setLogoutLoading(false);
    }
  };

  return (
    <Styled.Container>
      <Tooltip title="닫기" arrow disableInteractive>
        <IconButton
          onClick={() => {
            setMyInfoModal(false);
          }}
          aria-label="delete"
          sx={{
            position: "absolute",
            top: ".2rem",
            right: ".4rem",

            width: "12px",
            height: "12px",
          }}
        >
          <CloseIcon sx={{ fontSize: 16 }} />
        </IconButton>
      </Tooltip>
      <Styled.InfoTop>
        <Styled.ProfileImgBox>
          <img src="/images/logo.webp" alt="profile" />
        </Styled.ProfileImgBox>
        <Styled.NameContainer>
          <div style={{ display: "flex", alignItems: "center" }}>
            {userState.user.user.username}
            <div style={{ flexGrow: "1" }} />

            <Tooltip title="수정" arrow disableInteractive>
              <IconButton
                onClick={() => {
                  console.log(1);
                }}
                aria-label="delete"
                sx={{
                  color: "#333",
                  width: "20px",
                  height: "20px",
                }}
              >
                <EditIcon sx={{ fontSize: 14 }} />
              </IconButton>
            </Tooltip>
          </div>
          <div>{userState.user.user.email}</div>
        </Styled.NameContainer>
        <Styled.LogoutButtonContainer>
          <ActionButton bg="black" onClick={handleLogout}>
            {logoutLoading ? (
              <CircularProgress size={19.5} sx={{ color: "#fff" }} />
            ) : (
              "로그아웃"
            )}
          </ActionButton>
        </Styled.LogoutButtonContainer>
      </Styled.InfoTop>
      <Styled.InfoBottom>
        <Button sx={{ width: "33.33%", color: "#333" }}>주변 검색</Button>
        <Button sx={{ width: "33.33%", color: "#333" }}>내 장소</Button>
        <Button sx={{ width: "33.33%", color: "#333" }}>결제 정보</Button>
      </Styled.InfoBottom>
    </Styled.Container>
  );
};

export default MyInfoModal;
