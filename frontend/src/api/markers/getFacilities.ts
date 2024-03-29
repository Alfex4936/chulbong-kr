import instance from "../instance";
import type { Facilities } from "./setFacilities";

interface FacilitiesRes extends Facilities {
  markerId: number;
}

const getFacilities = async (markerId: number): Promise<FacilitiesRes[]> => {
  try {
    const res = await instance.get(`/api/v1/markers/${markerId}/facilities`, {
      withCredentials: true,
    });

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default getFacilities;
