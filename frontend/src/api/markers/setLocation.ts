import axios from "axios";

interface SetLocationReq {
  description: string;
  photoUrl: string;
  latitude: number;
  longitude: number;
}

export interface SetLocationRes {
  markerId: number;
  latitude: number;
  longitude: number;
  description: string;
  username: string;
  photoUrl: string;
}

const setLocation = async (body: SetLocationReq) => {
  const token = JSON.parse(localStorage.getItem("user") as string).state.user
    .token;

  try {
    const res = await axios.post(
      `${import.meta.env.VITE_LOCAL_URL}/api/v1/markers`,
      body,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    return res;
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      console.error(
        `위치 등록 실패: ${error.response.status} - ${error.response.data.error}`
      );

      return error.response.data.error;
    } else {
      console.error(`위치 등록 실패: ${error}`);
    }
  }
};

export default setLocation;
