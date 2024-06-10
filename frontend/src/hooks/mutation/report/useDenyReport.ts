import denyReport from "@/api/report/denyReport";
import { useMutation, useQueryClient } from "@tanstack/react-query";

const useDenyReport = (markerId: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: denyReport,
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

export default useDenyReport;
