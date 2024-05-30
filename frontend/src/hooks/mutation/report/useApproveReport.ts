import approveReport from "@/api/report/approveReport";
import useMapStore from "@/store/useMapStore";
import { useMutation, useQueryClient } from "@tanstack/react-query";

const useApproveReport = (markerId: number, lat: number, lng: number) => {
  const queryClient = useQueryClient();
  const { markers, overlay } = useMapStore();

  const filtering = () => {
    if (!markers || !overlay) return;
    const newPosition = new window.kakao.maps.LatLng(lat, lng);
    const marker = markers.find((value) => Number(value.Gb) === markerId);

    if (marker) {
      marker.setPosition(newPosition);
      overlay.setPosition(newPosition);
    }
  };

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
      filtering();
    },
  });
};

export default useApproveReport;
