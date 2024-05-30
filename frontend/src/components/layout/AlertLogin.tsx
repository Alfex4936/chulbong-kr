"use client";

import useLoginModalStateStore from "@/store/useLoginModalStateStore";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useRef } from "react";
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
import { Button } from "../ui/button";

const AlertLogin = () => {
  const router = useRouter();
  const pathname = usePathname();

  const { close, isOpen } = useLoginModalStateStore();

  const modalRef = useRef<HTMLButtonElement>(null);

  useEffect(() => {
    if (!modalRef) return;

    if (isOpen) modalRef.current?.click();
  }, [isOpen]);

  return (
    <AlertDialog>
      <AlertDialogTrigger asChild className="hidden">
        <Button variant="outline" ref={modalRef}>
          Show Dialog
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>로그인 완료 시 이용 가능합니다.</AlertDialogTitle>
          <AlertDialogDescription>
            로그인 후 여러 위치를 관리해보세요!
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel onClick={() => close()}>취소</AlertDialogCancel>
          <AlertDialogAction
            onClick={() => router.push(`/signin?redirect=${pathname}`)}
          >
            로그인 하러가기
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
};

export default AlertLogin;
