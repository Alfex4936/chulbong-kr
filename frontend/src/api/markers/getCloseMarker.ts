import instance from "../instance";

interface Props {
  lat: number;
  lng: number;
  distance: number;
  pageParam: number;
}

interface CloseMarker {
  latitude: number;
  longitude: number;
  distance: number;
  markerId: number;
  description: string;
  address?: string;
}

interface CloseMarkerRes {
  currentPage: number;
  markers: CloseMarker[];
  totalMarkers: number;
  totalPages: number;
}

const getCloseMarker = async ({
  lat,
  lng,
  distance,
  pageParam,
}: Props): Promise<CloseMarkerRes> => {
  const res = await instance.get(
    `/api/v1/markers/close?latitude=${lat}&longitude=${lng}&distance=${distance}&n=5&page=${pageParam}&pageSize=10`
  );

  return res.data;
};

export default getCloseMarker;
