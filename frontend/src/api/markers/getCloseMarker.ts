import instance from "../instance";

interface Props {
  lat: number;
  lon: number;
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
  lon,
  distance,
  pageParam,
}: Props): Promise<CloseMarkerRes> => {
  try {
    const res = await instance.get(
      `/api/v1/markers/close?latitude=${lat}&longitude=${lon}&distance=${distance}&n=5&page=${pageParam}`
    );

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default getCloseMarker;
