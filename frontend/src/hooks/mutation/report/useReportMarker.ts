import reportMarker from "@/api/report/reportMarker";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";

const useReportMarker = () => {
  const router = useRouter();
  return useMutation({
    mutationFn: reportMarker,
    onSuccess: () => {
      router.push("/mypage/report");
    },
  });
};

export default useReportMarker;
