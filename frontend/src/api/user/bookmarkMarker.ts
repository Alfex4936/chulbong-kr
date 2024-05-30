import instance from "../instance";

export interface Favorite {
  latitude: number;
  longitude: number;
  markerId: number;
  description: string;
  address?: string;
}

const bookmarkMarker = async (): Promise<Favorite[]> => {
  const res = await instance.get(`/api/v1/users/favorites`);

  return res.data;
};

export default bookmarkMarker;
