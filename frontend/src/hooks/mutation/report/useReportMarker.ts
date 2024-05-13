import reportMarker from "@/api/report/reportMarker";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";

const useReportMarker = () => {
  const router = useRouter();
  return useMutation({
    mutationFn: reportMarker,
    onSuccess: () => {
      // TODO: 마커 상세보기 페이지 리포트 리스트로 변경
      // TODO: 리포트  query invalidate
      router.push("/mypage/report");
    },
  });
};

export default useReportMarker;
