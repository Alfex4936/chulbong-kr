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
    const res = await axios.post(
      `${import.meta.env.VITE_LOCAL_URL}/api/v1/auth/signup`,
      body
    );

    return { data: res };
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      console.error(
        `회원가입 실패: ${error.response.status} - ${error.response.data.error}`
      );

      return {
        error: {
          code: error.response.status,
          msg: error.response.data.error,
        },
      };
    } else {
      console.error(`회원가입 실패: ${error}`);
      return { error: { code: 0, msg: "Unknown error" } };
    }
  }
};

export default signin;
