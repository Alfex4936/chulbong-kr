import deleteMarker from "@/api/markers/deleteMarker";
import { useToast } from "@/components/ui/use-toast";
import useLoginModalStateStore from "@/store/useLoginModalStateStore";
import useMapStore from "@/store/useMapStore";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { isAxiosError } from "axios";
import { useRouter } from "next/navigation";

interface Marker {
  address: string;
  description: string;
  latitude: number;
  longitude: number;
  markerId: number;
}

interface Page {
  currentPage: number;
  markers: Marker[];
  totalMarkers: number;
  totalPages: number;
}

interface Data {
  pages: Page[];
  pageParams: number[];
}

const useDeleteMarker = ({
  id,
  isRouting = false,
}: {
  id: number;
  isRouting?: boolean;
}) => {
  const router = useRouter();

  const { open } = useLoginModalStateStore();
  const { toast } = useToast();
  const { clusterer, markers, setMarkers, overlay, setOverlay } = useMapStore();

  const queryClient = useQueryClient();

  const filtering = async () => {
    if (!markers || !clusterer) return;
    const marker = markers.find((value) => Number(value.Gb) === id);

    const newMarkers = markers.filter((value) => Number(value.Gb) !== id);

    if (marker) {
      marker.setMap(null);
      clusterer.removeMarker(marker);
      if (overlay) {
        overlay.setMap(null);
      }
      setMarkers(newMarkers);
      setOverlay(null);
    }
  };

  return useMutation({
    mutationFn: () => {
      return deleteMarker(id);
    },

    onMutate: async () => {
      await queryClient.cancelQueries({ queryKey: ["myMarker"] });

      const previousMarkerData: Data = queryClient.getQueryData([
        "myMarker",
      ]) as Data;

      if (!previousMarkerData) return;

      const copy = { ...previousMarkerData };
      copy.pages.forEach((page) => {
        page.markers = page.markers.filter((marker) => marker.markerId !== id);
      });

      queryClient.setQueryData(["myMarker"], copy);

      return { previousMarkerData };
    },

    onError: (error, _hero, context?: { previousMarkerData: Data }) => {
      console.log(error);
      if (isAxiosError(error)) {
        if (error.response?.status === 401) {
          open();
        } else {
          toast({ description: "잠시 후 다시 시도해 주세요." });
        }
      } else {
        toast({ description: "잠시 후 다시 시도해 주세요." });
      }
      if (context?.previousMarkerData) {
        queryClient.setQueryData(["myMarker"], context.previousMarkerData);
      }
    },

    onSuccess: async () => {
      await filtering();
      if (isRouting) router.replace("/home");
      toast({ description: "삭제가 완료됐습니다." });
    },

    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["myMarker"] });
    },
  });
};

export default useDeleteMarker;
