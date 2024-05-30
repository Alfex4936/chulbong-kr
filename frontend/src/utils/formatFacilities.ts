import { type FacilitiesRes } from "@/api/markers/getFacilities";

export interface FormatedFacilities {
  철봉: number | string;
  평행봉: number | string;
}

const formatFacilities = (facilities: FacilitiesRes[]): FormatedFacilities => {
  const data: { 철봉: number | string; 평행봉: number | string } = {
    철봉: 0,
    평행봉: 0,
  };

  data.철봉 =
    (facilities?.find((facilitie) => facilitie.facilityId === 1)
      ?.quantity as number) || "개수 정보 없음";
  data.평행봉 =
    (facilities?.find((facilitie) => facilitie.facilityId === 2)
      ?.quantity as number) || "개수 정보 없음";

  return data;
};

export default formatFacilities;
