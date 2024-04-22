import instance from "../instance";
import type { Facilities } from "./setFacilities";

export interface FacilitiesRes extends Facilities {
  markerId: number;
}

const getFacilities = async (markerId: number): Promise<FacilitiesRes[]> => {
  const res = await instance.get(`/api/v1/markers/${markerId}/facilities`, {
    withCredentials: true,
  });

  return res.data;
};

export default getFacilities;
