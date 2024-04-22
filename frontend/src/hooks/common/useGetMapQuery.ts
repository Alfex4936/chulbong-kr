import { useSearchParams } from "next/navigation";

const useGetMapQuery = () => {
  const searchParams = useSearchParams();

  const lat = searchParams.get("lat") || 37.566535;
  const lng = searchParams.get("lng") || 126.9779692;
  const level = searchParams.get("lv") || 3;

  return { lat, lng, level };
};

export default useGetMapQuery;
