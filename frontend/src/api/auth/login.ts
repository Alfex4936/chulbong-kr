import axios from "axios";

interface LoginReq {
  email: string;
  password: string;
}

export interface User {
  userId: number;
  username: string;
  email: string;
}

export interface LoginRes {
  token: string;
  user: User;
}

const login = async (body: LoginReq) => {
  try {
    const res = await axios.post(
      `${import.meta.env.VITE_LOCAL_URL}/api/v1/auth/login`,
      body
    );

    return res;
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      console.error(
        `로그인 실패: ${error.response.status} - ${error.response.data.error}`
      );

      return error.response.data.error;
    } else {
      console.error(`로그인 실패: ${error}`);
    }
  }
};

export default login;
