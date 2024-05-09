import setNewMarker from "@/api/markers/setNewMarker";
import { useToast } from "@/components/ui/use-toast";
import useLoginModalStateStore from "@/store/useLoginModalStateStore";
import { useMutation } from "@tanstack/react-query";
import { isAxiosError } from "axios";

const useUploadMarker = () => {
  const { toast } = useToast();
  const { open } = useLoginModalStateStore();

  return useMutation({
    mutationFn: setNewMarker,
    onError: (error) => {
      console.log(error);
      if (isAxiosError(error)) {
        if (error.response?.status === 401) {
          open();
        } else {
          toast({ description: "잠시 후 다시 시도해 주세요." });
        }
      } else {
        toast({ description: "잠시 후 다시 시도해 주세요." });
      }
    },
  });
};

export default useUploadMarker;
