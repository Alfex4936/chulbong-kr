import denyReport from "@/api/report/denyReport";
import { useMutation } from "@tanstack/react-query";

const useDenyReport = () => {
  return useMutation({
    mutationFn: denyReport,
  });
};

export default useDenyReport;
