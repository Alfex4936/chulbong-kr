import deleteReport from "@/api/report/deleteReport";
import { MutationFunction, useMutation } from "@tanstack/react-query";

const useDeleteReport = () => {
  return useMutation({
    mutationFn: deleteReport as MutationFunction<any, number>,
  });
};

export default useDeleteReport;
