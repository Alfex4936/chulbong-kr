import CloseIcon from "@mui/icons-material/Close";
import Button from "@mui/material/Button";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { useState } from "react";
import useGetMyInfo from "../../hooks/query/user/useGetMyInfo";
import type { KakaoMap } from "../../types/KakaoMap.types";
import AroundMarker from "../AroundMarker/AroundMarker";
import MyInfoDetail from "../MyInfoDetail/MyInfoDetail";
import MyMarker from "../MyMarker/MyMarker";
import * as Styled from "./MyInfoModal.style";

interface Props {
  map: KakaoMap;
  setMyInfoModal: React.Dispatch<React.SetStateAction<boolean>>;
  setDeleteUserModal: React.Dispatch<React.SetStateAction<boolean>>;
}

const MyInfoModal = ({ map, setMyInfoModal, setDeleteUserModal }: Props) => {
  const { data, isLoading } = useGetMyInfo();

  const tabs = [
    { title: "주변 검색", content: <AroundMarker map={map} /> },
    { title: "내 장소", content: <MyMarker map={map} /> },
    {
      title: "내 정보",
      content: (
        <MyInfoDetail
          setDeleteUserModal={setDeleteUserModal}
          setMyInfoModal={setMyInfoModal}
        />
      ),
    },
  ];

  const [curTab, setCurTab] = useState<number | null>(null);

  return (
    <Styled.Container>
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
        <Styled.ButtonContainer>
          <Tooltip title="닫기" arrow disableInteractive>
            <IconButton
              onClick={() => {
                setMyInfoModal(false);
              }}
              aria-label="delete"
            >
              <CloseIcon />
            </IconButton>
          </Tooltip>
        </Styled.ButtonContainer>
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
