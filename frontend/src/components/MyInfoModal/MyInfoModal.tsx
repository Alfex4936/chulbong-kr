import CloseIcon from "@mui/icons-material/Close";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { useState } from "react";
import logout from "../../api/auth/logout";
import useUserStore from "../../store/useUserStore";
import type { KakaoMap } from "../../types/KakaoMap.types";
import ActionButton from "../ActionButton/ActionButton";
import AroundMarker from "../AroundMarker/AroundMarker";
import MyInfoDetail from "../MyInfoDetail/MyInfoDetail";
import MyMarker from "../MyMarker/MyMarker";
import * as Styled from "./MyInfoModal.style";
import useGetMyInfo from "../../hooks/query/user/useGetMyInfo";

interface Props {
  map: KakaoMap;
  setMyInfoModal: React.Dispatch<React.SetStateAction<boolean>>;
  setDeleteUserModal: React.Dispatch<React.SetStateAction<boolean>>;
}

const MyInfoModal = ({ map, setMyInfoModal, setDeleteUserModal }: Props) => {
  const userState = useUserStore();
  const { data, isLoading } = useGetMyInfo();

  const tabs = [
    { title: "주변 검색", content: <AroundMarker map={map} /> },
    { title: "내 장소", content: <MyMarker map={map} /> },
    {
      title: "내 정보",
      content: <MyInfoDetail setDeleteUserModal={setDeleteUserModal} />,
    },
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
        {isLoading ? (
          <Styled.NameSkeleton>
            <div />
            <div />
          </Styled.NameSkeleton>
        ) : (
          <Styled.NameContainer>
            <div style={{ display: "flex", alignItems: "center" }}>
              {data?.username}
              <div style={{ flexGrow: "1" }} />
            </div>
            <div>{data?.email}</div>
          </Styled.NameContainer>
        )}
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