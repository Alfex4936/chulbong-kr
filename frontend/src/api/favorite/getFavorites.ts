import axios from "axios";

interface Favorite {
  latitude: number;
  longitude: number;
  markerId: number;
  description: string;
  address?: string;
}

const getFavorites = async (): Promise<Favorite[]> => {
  try {
    const res = await axios.get(`/api/v1/users/favorites`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default getFavorites;
