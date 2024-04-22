import { type FacilitiesRes } from "@/api/markers/getFacilities";

export interface FormatedFacilities {
  철봉: number;
  평행봉: number;
}

const formatFacilities = (facilities: FacilitiesRes[]): FormatedFacilities => {
  const data = { 철봉: 0, 평행봉: 0 };
  data.철봉 = facilities.find((facilitie) => facilitie.facilityId === 1)
    ?.quantity as number;
  data.평행봉 = facilities.find((facilitie) => facilitie.facilityId === 2)
    ?.quantity as number;

  return data;
};

export default formatFacilities;
