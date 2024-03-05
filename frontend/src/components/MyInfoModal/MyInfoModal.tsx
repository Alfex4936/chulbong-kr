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
import AroundMarker from "../AroundMarker/AroundMarker";
import MyMarker from "../MyMarker/MyMarker";
import PaymentInfo from "../PaymentInfo/PaymentInfo";
import * as Styled from "./MyInfoModal.style";
import type { KakaoMap } from "../../types/KakaoMap.types";

interface Props {
  map: KakaoMap;
  setMyInfoModal: React.Dispatch<React.SetStateAction<boolean>>;
}

const MyInfoModal = ({ map, setMyInfoModal }: Props) => {
  const userState = useUserStore();
  const tabs = [
    { title: "주변 검색", content: <AroundMarker /> },
    { title: "내 장소", content: <MyMarker map={map} /> },
    { title: "결제 정보", content: <PaymentInfo /> },
  ];

  const [logoutLoading, setLogoutLoading] = useState(false);
  const [curTab, setCurTab] = useState<number | null>(null);

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
        {tabs.map((tab, index) => {
          return (
            <Button
              key={index}
              sx={{
                width: "33.33%",
                color: index === curTab ? "#6b73db" : "#333",
              }}
              onClick={() => {
                setCurTab(index);
              }}
            >
              {tab.title}
            </Button>
          );
        })}
      </Styled.InfoBottom>
      {curTab !== null && (
        <Styled.TabContainer>{tabs[curTab].content}</Styled.TabContainer>
      )}
    </Styled.Container>
  );
};

export default MyInfoModal;
