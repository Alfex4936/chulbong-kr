import EditIcon from "@mui/icons-material/Edit";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { useState } from "react";
import useUpdateName from "../../hooks/mutation/user/useUpdateName";
import useGetMyInfo from "../../hooks/query/user/useGetMyInfo";
import useInput from "../../hooks/useInput";
import ActionButton from "../ActionButton/ActionButton";
import Input from "../Input/Input";
import * as Styled from "./MyInfoDetail.style";

interface Props {
  setDeleteUserModal: React.Dispatch<React.SetStateAction<boolean>>;
}

const MyInfoDetail = ({ setDeleteUserModal }: Props) => {
  const nameValue = useInput("");

  const { data, isLoading } = useGetMyInfo();
  const { mutate, isPending } = useUpdateName(nameValue.value);

  const [updateName, setUpdateName] = useState(false);
  const [updateNameError, setUpdateNameError] = useState("");

  const handleUpdateName = () => {
    if (nameValue.value === data?.username) {
      setUpdateNameError("현재 닉네임과 동일합니다.");
    } else {
      mutate();
      setUpdateName(false);
    }
  };

  const handleAlert = () => {
    console.log(2);
    setDeleteUserModal(true);
    // Swal.fire({
    //   title: "정말 삭제하시겠습니까?",
    //   text: "되돌릴 수 없습니다!",
    //   icon: "warning",
    //   showCancelButton: true,
    //   confirmButtonColor: "#3085d6",
    //   cancelButtonColor: "#d33",
    //   confirmButtonText: "삭제",
    //   cancelButtonText: "취소",
    // }).then((result) => {
    //   if (result.isConfirmed) {
    //     Swal.fire({
    //       title: "삭제 완료",
    //       icon: "success",
    //     });
    //   }
    // });
  };

  const handleUpdatePassword = () => {
    console.log(1);
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
                <Styled.Name>{data?.username}</Styled.Name>
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
        <ActionButton bg="black" onClick={handleUpdatePassword}>
          비밀번호 변경
        </ActionButton>
        <ActionButton bg="gray" onClick={handleAlert}>
          회원 탈퇴
        </ActionButton>
      </Styled.ButtonContainer>
    </Styled.Container>
  );
};

export default MyInfoDetail;
