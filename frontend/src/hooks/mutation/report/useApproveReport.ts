import approveReport from "@/api/report/approveReport";
import { useMutation, useQueryClient } from "@tanstack/react-query";

const useApproveReport = (markerId: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: approveReport,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["marker", "report", "me"] });
      queryClient.invalidateQueries({ queryKey: ["marker", "report", "all"] });
      queryClient.invalidateQueries({
        queryKey: ["marker", "report", "formarker"],
      });
      queryClient.invalidateQueries({
        queryKey: ["marker", "report", "formarker", markerId],
      });
      queryClient.invalidateQueries({
        queryKey: ["marker", markerId],
      });
      queryClient.invalidateQueries({
        queryKey: ["markers"],
      });
      queryClient.invalidateQueries({
        queryKey: ["report", "my"],
      });
    },
  });
};

export default useApproveReport;
