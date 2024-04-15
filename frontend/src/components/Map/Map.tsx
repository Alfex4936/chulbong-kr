import getAllMarker from "@/api/markers/getAllMarker";
import {
  HydrationBoundary,
  QueryClient,
  dehydrate,
} from "@tanstack/react-query";
import MapClient from "./MapClient";

const Map = async () => {
  const queryClient = new QueryClient();
  await queryClient.prefetchQuery({
    queryKey: ["markers"],
    queryFn: getAllMarker,
  });

  const dehydrateState = dehydrate(queryClient);

  return (
    <HydrationBoundary state={dehydrateState}>
      <MapClient />
    </HydrationBoundary>
  );
};

export default Map;
