import EditIcon from "@/components/icons/EditIcon";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import useInput from "@/hooks/common/useInput";
import { useEffect, useState } from "react";
import useSendPasswordReset from "@/hooks/mutation/auth/useSendPasswordReset";
import LoadingSpinner from "@/components/atom/LoadingSpinner";

interface Props {
  text?: string;
  textClass?: string;
}

const ChangePassword = ({ text, textClass }: Props) => {
  const { mutateAsync: sendEmail, isPending: isSending } =
    useSendPasswordReset();
  const { value, handleChange, resetValue } = useInput("");

  const [isSended, setIsSended] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");

  useEffect(() => {
    setIsSended(false);
  }, []);

  useEffect(() => {
    if (errorMessage === "") return;

    const timeout = setTimeout(() => {
      setErrorMessage("");
    }, 3000);

    return () => clearTimeout(timeout);
  }, [errorMessage]);

  const handleSendEmail = async () => {
    try {
      await sendEmail(value);
      setIsSended(true);
      resetValue();
    } catch (error) {
      setErrorMessage("잠시 후 다시 시도해 주세요.");
    }
  };

  return (
    <Dialog>
      <DialogTrigger asChild>
        {text && textClass ? (
          <button className={textClass}>{text}</button>
        ) : (
          <button className="p-1 rounded-full hover:bg-white-tp-light">
            <EditIcon size={12} />
          </button>
        )}
      </DialogTrigger>

      {isSended ? (
        <DialogContent className="w-3/4 min-w-80 web:max-w-[425px] bg-black-dark">
          <DialogHeader>
            <DialogTitle>전송 완료</DialogTitle>
            <DialogDescription>이메일을 확인해 주세요.</DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              className=""
              variant="secondary"
              onClick={() => setIsSended(false)}
            >
              이전
            </Button>
            <DialogClose asChild>
              <Button className="">닫기</Button>
            </DialogClose>
          </DialogFooter>
        </DialogContent>
      ) : (
        <DialogContent className="w-3/4 min-w-80 web:max-w-[425px] bg-black-dark">
          <DialogHeader>
            <DialogTitle>비밀번호 초기화</DialogTitle>
            <DialogDescription>
              이메일로 비밀번호 변경 링크가 전송됩니다. <br />
              초기화하고자 하는 이메일을 입력해 주세요.
            </DialogDescription>
          </DialogHeader>
          <div className="flex-col">
            <div className="flex items-center">
              <Label htmlFor="email" className="text-left w-1/5">
                이메일
              </Label>
              <Input
                type="email"
                value={value}
                onChange={handleChange}
                id="email"
                className="w-4/5"
              />
            </div>
            <div className="text-sm text-red mt-1">{errorMessage}</div>
          </div>
          <DialogFooter>
            <Button onClick={handleSendEmail}>
              {isSending ? (
                <LoadingSpinner size="xs" color="black" />
              ) : (
                "메일 보내기"
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      )}
    </Dialog>
  );
};

export default ChangePassword;
