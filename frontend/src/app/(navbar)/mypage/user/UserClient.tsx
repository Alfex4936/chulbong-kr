import BlackLightBox from "@/components/atom/BlackLightBox";
import EmojiHoverButton from "@/components/atom/EmojiHoverButton";
import GrowBox from "@/components/atom/GrowBox";
import EditIcon from "@/components/icons/EditIcon";
import { Separator } from "@/components/ui/separator";

interface Props {
  text: string;
  subText: string;
  buttonText?: string;
  onClick?: VoidFunction;
}

const InfoList = ({ text, subText, buttonText, onClick }: Props) => {
  return (
    <div className="flex text-[13px] py-1">
      <span className="w-16">{text}</span>
      <span className="">{subText}</span>
      <GrowBox />
      {buttonText && (
        <button className="p-1 rounded-full hover:bg-white-tp-light">
          <EditIcon size={12} />
        </button>
      )}
    </div>
  );
};

const UserClient = () => {
  return (
    <div>
      <div className="mb-4 mt-5">
        <BlackLightBox center>
          <div className="flex justify-center items-center text-xl mb-2">
            <span className="mr-2">이용훈</span>
            <button className="p-1 rounded-full hover:bg-white-tp-light">
              <EditIcon size={15} />
            </button>
          </div>
          <div className="text-sm text-grey-dark">yonghuni484@gmail.com</div>
        </BlackLightBox>
      </div>

      <div className="mb-4">
        <BlackLightBox>
          <div>개인 정보</div>
          <Separator className="mx-1 my-3 bg-grey-dark-1" />
          <InfoList text="아이디" subText="yonghuni484@gmail.com" />
          <InfoList text="이메일" subText="yonghuni484@gmail.com" />
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
