import instance from "../instance";

export interface Facilities {
  facilityId: number;
  quantity: number;
}

const setFacilities = async (markerId: number, facilities: Facilities[]) => {
  const body = {
    markerId,
    facilities,
  };

  const res = await instance.post(`/api/v1/markers/facilities`, body, {
    withCredentials: true,
  });

  return res.data;
};

export default setFacilities;
