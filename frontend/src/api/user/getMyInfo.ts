import axios from "axios";

interface MyInfo {
  userId: number;
  username: string;
  email: string;
}

const getMyInfo = async (): Promise<MyInfo> => {
  try {
    const res = await axios.get(`/api/v1/users/me`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default getMyInfo;
