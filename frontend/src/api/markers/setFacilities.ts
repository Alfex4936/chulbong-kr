import axios from "axios";

export interface Facilities {
  facilityId: number;
  quantity: number;
}

const setFacilities = async (markerId: number, facilities: Facilities[]) => {
  const body = {
    markerId,
    facilities,
  };

  try {
    const res = await axios.post(`/api/v1/markers/facilities`, body, {
      withCredentials: true,
    });

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default setFacilities;
