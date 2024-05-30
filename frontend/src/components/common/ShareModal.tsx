import { cn } from "@/lib/utils";
import useLoginModalStateStore from "@/store/useLoginModalStateStore";
import axios, { isAxiosError } from "axios";
import { RefObject, useEffect, useRef, useState } from "react";
import GrowBox from "../atom/GrowBox";
import LoadingSpinner from "../atom/LoadingSpinner";
import ExitIcon from "../icons/ExitIcon";
import { Button } from "../ui/button";
import { useToast } from "../ui/use-toast";

interface Props {
  link: string;
  className?: string;
  buttonRef: RefObject<HTMLDivElement>;
  lat: number;
  lng: number;
  filename: string;
  closeModal: VoidFunction;
}

const ShareModal = ({
  link,
  className,
  buttonRef,
  lat,
  lng,
  filename,
  closeModal,
}: Props) => {
  const { open: openLoginModal } = useLoginModalStateStore();
  const { toast } = useToast();

  const modalRef = useRef<HTMLDivElement>(null);

  const [downLoading, setDownLoading] = useState(false);

  useEffect(() => {
    const handleClick = (e: MouseEvent) => {
      if (
        modalRef.current &&
        buttonRef.current &&
        !modalRef.current.contains(e.target as Node) &&
        !buttonRef.current.contains(e.target as Node)
      ) {
        closeModal();
      }
    };
    window.addEventListener("click", handleClick);

    return () => {
      window.removeEventListener("click", handleClick);
    };
  }, [modalRef, buttonRef]);

  const copyTextToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(link);
      toast({
        description: "링크 복사 완료",
      });
    } catch (err) {
      alert("잠시 후 다시 시도해 주세요!");
    }
  };

  const downloadMap = async () => {
    setDownLoading(true);
    let downlink;
    let link;
    try {
      const response = await axios.get(
        `/api/v1/markers/save-offline?latitude=${lat}&longitude=${lng}`,
        {
          responseType: "blob",
        }
      );

      downlink = window.URL.createObjectURL(new Blob([response.data]));
      link = document.createElement("a");

      link.href = downlink;
      link.download = `${filename}.pdf`;

      document.body.appendChild(link);

      link.click();
    } catch (error) {
      if (isAxiosError(error)) {
        if (error.response?.status === 401) {
          openLoginModal();
        } else {
          toast({ description: "잠시 후 다시 시도해 주세요." });
        }
      } else {
        toast({ description: "잠시 후 다시 시도해 주세요." });
      }
    } finally {
      if (link) document.body.removeChild(link);
      if (downlink) window.URL.revokeObjectURL(downlink);
      setDownLoading(false);
    }
  };

  return (
    <div
      className={cn("bg-grey text-black p-2 rounded-sm shadow-md", className)}
      ref={modalRef}
    >
      <div className="flex">
        <Button
          className="flex items-center justify-center w-[87px] h-[37] mr-3 hover:bg-grey-dark"
          onClick={copyTextToClipboard}
        >
          링크복사
        </Button>
        <Button
          className="flex items-center justify-center w-[87px] h-[37] hover:bg-grey-dark"
          onClick={downloadMap}
          disabled={downLoading}
        >
          {downLoading ? (
            <LoadingSpinner size="xs" color="black" />
          ) : (
            "PDF 다운"
          )}
        </Button>
        <GrowBox />
        <button className="flex items-start" onClick={closeModal}>
          <ExitIcon size={18} color="black" />
        </button>
      </div>
      <p className="text-xs text-red mt-1">
        해당 위치를 기반으로 주변의 지도와 철봉 위치들을
        <br />
        PDF로 다운로들 할 수 있습니다.
      </p>
    </div>
  );
};

export default ShareModal;
