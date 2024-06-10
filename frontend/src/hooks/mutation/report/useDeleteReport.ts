import deleteReport from "@/api/report/deleteReport";
import { useMutation, useQueryClient } from "@tanstack/react-query";

const useDeleteReport = (markerId: number, reportId: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      return deleteReport(markerId, reportId);
    },
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
        queryKey: ["report", "my"],
      });
    },
  });
};

export default useDeleteReport;
