import reportMarker from "@/api/report/reportMarker";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";

const useReportMarker = (markerId: number) => {
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: reportMarker,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["marker", "report", "me"] });
      queryClient.invalidateQueries({ queryKey: ["marker", "report", "all"] });
      queryClient.invalidateQueries({
        queryKey: ["marker", "report", "formarker"],
      });
      queryClient.invalidateQueries({
        queryKey: ["marker", "report", "formarker", markerId],
      });

      router.push(`/pullup/${markerId}/reportlist`);
    },
  });
};

export default useReportMarker;
