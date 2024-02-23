import axios, { AxiosResponse } from "axios";

interface SigninReq {
  username?: string;
  email: string;
  password: string;
}

export interface SigninRes {
  userId: number;
  username: string;
  email: string;
}

interface SigninResponse {
  data?: AxiosResponse<SigninRes>;
  error?: { code: number; msg: string };
}

const signin = async (body: SigninReq): Promise<SigninResponse> => {
  try {
    const res = await axios.post(`/api/v1/auth/signup`, body);

    return { data: res };
  } catch (error) {
    throw error;
  }
};

export default signin;
