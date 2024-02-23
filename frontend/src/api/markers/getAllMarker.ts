import axios from "axios";

const getAllMarker = async () => {
  try {
    const res = await axios.get(
      `${import.meta.env.VITE_LOCAL_URL}/api/v1/markers`
    );

    return res;
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      console.error(
        `마커 불러오기 실패: ${error.response.status} - ${error.response.data.error}`
      );
    } else {
      console.error(`마커 불러오기 실패: ${error}`);
    }
  }
};

export default getAllMarker;
