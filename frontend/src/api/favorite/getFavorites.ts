import instance from "../instance";

interface Favorite {
  latitude: number;
  longitude: number;
  markerId: number;
  description: string;
  address?: string;
}

const getFavorites = async (): Promise<Favorite[]> => {
  try {
    const res = await instance.get(`/api/v1/users/favorites`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default getFavorites;
