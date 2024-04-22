export interface AddressInfo {
  address_name: string;
  region_1depth_name: string;
  region_2depth_name: string;
  region_3depth_name: string;
  mountain_yn: "Y" | "N";
  main_address_no: string;
  sub_address_no: string;
  zip_code: string;
}

const getAddress = (
  lat: number,
  lng: number
): Promise<AddressInfo | string> => {
  return new Promise((resolve) => {
    let geocoder = new window.kakao.maps.services.Geocoder();
    let coord = new window.kakao.maps.LatLng(lat, lng);

    geocoder.coord2Address(
      coord.getLng(),
      coord.getLat(),
      (result: { address: AddressInfo }[], status: string) => {
        if (status === window.kakao.maps.services.Status.OK) {
          resolve(result[0].address);
        } else {
          resolve("주소 정보 없음");
        }
      }
    );
  });
};

export default getAddress;
