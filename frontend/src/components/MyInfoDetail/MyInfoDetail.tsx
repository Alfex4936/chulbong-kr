import EditIcon from "@mui/icons-material/Edit";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { useState } from "react";
import useLogout from "../../hooks/mutation/auth/useLogout";
import useUpdateName from "../../hooks/mutation/user/useUpdateName";
import useGetMyInfo from "../../hooks/query/user/useGetMyInfo";
import useInput from "../../hooks/useInput";
import useModalStore from "../../store/useModalStore";
import useUserStore from "../../store/useUserStore";
import ActionButton from "../ActionButton/ActionButton";
import Input from "../Input/Input";
import * as Styled from "./MyInfoDetail.style";

interface Props {
  setDeleteUserModal: React.Dispatch<React.SetStateAction<boolean>>;
  setMyInfoModal: React.Dispatch<React.SetStateAction<boolean>>;
}

const MyInfoDetail = ({ setDeleteUserModal, setMyInfoModal }: Props) => {
  const nameValue = useInput("");
  const modalState = useModalStore();
  const userState = useUserStore();

  const { data, isLoading } = useGetMyInfo();
  const { mutate, isPending } = useUpdateName(nameValue.value);
  const { mutateAsync: logout, isPending: logoutPending } = useLogout();

  const [updateName, setUpdateName] = useState(false);
  const [updateNameError, setUpdateNameError] = useState("");

  // const [logoutLoading, setLogoutLoading] = useState(false);

  const handleLogout = async () => {
    try {
      await logout();
      userState.resetUser();
      setMyInfoModal(false);
    } catch (error) {
      userState.resetUser();
      setMyInfoModal(false);
    }
  };

  const handleUpdateName = () => {
    if (nameValue.value === data?.username) {
      setUpdateNameError("현재 닉네임과 동일합니다.");
    } else {
      mutate();
      setUpdateName(false);
    }
  };

  const handleAlert = () => {
    setDeleteUserModal(true);
  };

  const handleUpdatePassword = () => {
    modalState.openPassword();
  };

  if (isLoading) {
    return (
      <Styled.ListSkeleton>
        <div />
        <div />
      </Styled.ListSkeleton>
    );
  }

  return (
    <Styled.Container>
      <Styled.NameContainer>
        {updateName ? (
          <div style={{ margin: "auto" }}>
            <Input
              type="text"
              defaultValue={data?.username}
              id="text"
              onChange={nameValue.onChange}
              style={{
                width: "80%",
              }}
            />
            <Styled.ErrorBox>{updateNameError}</Styled.ErrorBox>
            <Styled.NameButtonContainer>
              <ActionButton
                bg="black"
                disabled={nameValue.value === ""}
                onClick={handleUpdateName}
              >
                확인
              </ActionButton>
              <ActionButton
                bg="gray"
                onClick={() => {
                  setUpdateName(false);
                }}
              >
                취소
              </ActionButton>
            </Styled.NameButtonContainer>
          </div>
        ) : (
          <>
            {isPending ? (
              <Styled.ListSkeleton>
                <div />
                <div />
              </Styled.ListSkeleton>
            ) : (
              <>
                <Styled.Name>
                  <span>닉네임</span> : {data?.username}
                </Styled.Name>
                <Tooltip title="수정" arrow disableInteractive>
                  <IconButton
                    onClick={() => {
                      setUpdateName(true);
                    }}
                    aria-label="delete"
                    sx={{
                      color: "#333",
                      width: "30px",
                      height: "30px",
                    }}
                  >
                    <EditIcon sx={{ fontSize: 18 }} />
                  </IconButton>
                </Tooltip>
              </>
            )}
          </>
        )}
      </Styled.NameContainer>
      <Styled.EmailContainer>
        <div>이메일</div>
        <div>{data?.email}</div>
      </Styled.EmailContainer>

      {/* <Styled.PaymentContainer>
        <div>결제 정보</div>
        <div>.</div>
        <div>.</div>
        <div>준비중</div>
      </Styled.PaymentContainer> */}

      <Styled.ButtonContainer>
        <Styled.ButtonTop>
          <ActionButton bg="black" onClick={handleUpdatePassword}>
            비밀번호 변경
          </ActionButton>
          <ActionButton bg="gray" onClick={handleAlert}>
            회원 탈퇴
          </ActionButton>
        </Styled.ButtonTop>
        <Styled.ButtonBottom>
          <Button
            onClick={handleLogout}
            sx={{
              color: "#333",
              width: "100%",
              border: "1px solid #ccc",
              fontSize: ".8rem",
            }}
          >
            {logoutPending ? (
              <CircularProgress size={19.5} sx={{ color: "#333" }} />
            ) : (
              "로그아웃"
            )}
          </Button>
        </Styled.ButtonBottom>
      </Styled.ButtonContainer>

      <Styled.InfoContainer>
        <p>© 2024 chulbong-kr. All rights reserved.</p>
        <p>문의: chulbong.kr@gmail.com</p>
      </Styled.InfoContainer>
    </Styled.Container>
  );
};

export default MyInfoDetail;
