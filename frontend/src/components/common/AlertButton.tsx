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
} from "../ui/alert-dialog";

interface Props {
  ButtonText: string | React.ReactNode;
  title: string;
  desc?: string;
  approveText?: string;
  cancelText?: string;
  clickFn?: VoidFunction;
  disabled?: boolean;
}

const AlertButton = ({
  ButtonText,
  desc,
  title,
  approveText,
  cancelText,
  clickFn,
  disabled = false,
}: Props) => {
  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <button disabled={disabled}>{ButtonText}</button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{title}</AlertDialogTitle>
          <AlertDialogDescription className="text-red">
            {desc}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>{cancelText || "취소"}</AlertDialogCancel>
          <AlertDialogAction onClick={clickFn}>
            {approveText || "확인"}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
};

export default AlertButton;
