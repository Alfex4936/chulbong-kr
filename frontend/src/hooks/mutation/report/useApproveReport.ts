import approveReport from "@/api/report/approveReport";
import { useMutation } from "@tanstack/react-query";

const useApproveReport = () => {
  return useMutation({
    mutationFn: approveReport,
  });
};

export default useApproveReport;
