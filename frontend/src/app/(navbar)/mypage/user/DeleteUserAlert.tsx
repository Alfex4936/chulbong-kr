import EmojiHoverButton from "@/components/atom/EmojiHoverButton";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import useDeleteUser from "@/hooks/mutation/user/useDeleteUser";

const DeleteUserAlert = () => {
  const { mutate: deleteUser } = useDeleteUser();

  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <div className="w-[90%] mx-auto">
          <EmojiHoverButton
            emoji="❗"
            text="탈퇴하기"
            subText="다음에 만나요!"
            onClickFn={() => deleteUser()}
          />
        </div>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>정말 탈퇴하시겠습니까?</AlertDialogTitle>
          <AlertDialogDescription className="text-red">
            추가하신 마커는 유지되고, 작성한 댓글 밑 사진은 모두 삭제됩니다!
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>취소</AlertDialogCancel>
          <AlertDialogAction>탈퇴하기</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
};

export default DeleteUserAlert;
