"use client";

import BlackLightBox from "@/components/atom/BlackLightBox";
import EmojiHoverButton from "@/components/atom/EmojiHoverButton";
import GrowBox from "@/components/atom/GrowBox";
import LoadingSpinner from "@/components/atom/LoadingSpinner";
import EditIcon from "@/components/icons/EditIcon";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import useInput from "@/hooks/common/useInput";
import useUpdateUserName from "@/hooks/mutation/user/useUpdateUserName";
import useMyinfoData from "@/hooks/query/user/useMyinfoData";
import { useState } from "react";
import ChangePassword from "./ChangePassword";
import { useRouter } from "next/navigation";
import Link from "next/link";
// TODO: 에러 토스트 팝업 연결

interface Props {
  text: string;
  subText: string;
  buttonText?: string;
}

const InfoList = ({ text, subText, buttonText }: Props) => {
  return (
    <div className="flex text-[13px] py-1">
      <span className="w-16">{text}</span>
      <span className="">{subText}</span>
      <GrowBox />
      {buttonText && <ChangePassword />}
    </div>
  );
};

const UserClient = () => {
  const router = useRouter();
  const {
    mutate: updateName,
    isPending: isNameUpdate,
    isError: updateNameError,
    error,
  } = useUpdateUserName();
  const { data: myInfo, isError } = useMyinfoData();

  const { value, handleChange } = useInput(myInfo?.username as string);

  const [nameInput, setNameInput] = useState(false);

  if (updateNameError) console.log(error);

  if (isError) return <div>존재하지 않는 유저입니다.</div>;
  return (
    <div>
      <div className="mb-4 mt-5">
        <BlackLightBox center>
          <div className="flex justify-center items-center text-xl mb-2">
            {nameInput ? (
              <div>
                <Input
                  type="text"
                  className="border-grey mb-2"
                  value={value}
                  onChange={handleChange}
                  maxLength={8}
                />
                <div className="flex">
                  <Button
                    className="border-grey border bg-transparent dark:text-grey hover:bg-white-tp-light hover:border-transparent mr-2"
                    size={"sm"}
                    onClick={() => {
                      updateName(value);
                      setNameInput(false);
                    }}
                  >
                    변경
                  </Button>
                  <Button
                    className="border-grey border bg-transparent dark:text-grey hover:bg-white-tp-light hover:border-transparent"
                    size={"sm"}
                    onClick={() => setNameInput(false)}
                  >
                    취소
                  </Button>
                </div>
              </div>
            ) : (
              <>
                {isNameUpdate ? (
                  <LoadingSpinner size="sm" />
                ) : (
                  <>
                    <span className="mr-2">{myInfo?.username}</span>
                    <button
                      className="p-1 rounded-full hover:bg-white-tp-light"
                      onClick={() => setNameInput(true)}
                    >
                      <EditIcon size={15} />
                    </button>
                  </>
                )}
              </>
            )}
          </div>
          <div className="text-sm text-grey-dark">{myInfo?.email}</div>
        </BlackLightBox>
      </div>

      <div className="mb-4">
        <BlackLightBox>
          <div>개인 정보</div>
          <Separator className="mx-1 my-3 bg-grey-dark-1" />
          <InfoList text="아이디" subText={myInfo?.email as string} />
          <InfoList text="이메일" subText={myInfo?.email as string} />
          <InfoList
            text="비밀번호"
            subText="............"
            buttonText="수정하기"
          />
        </BlackLightBox>
      </div>

      <div className="w-[90%] mx-auto">
        <EmojiHoverButton emoji="❗" text="탈퇴하기" subText="다음에 만나요!" />
      </div>
    </div>
  );
};

export default UserClient;
