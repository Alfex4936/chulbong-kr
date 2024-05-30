import type { User } from "../../types/User.types";
import instance from "../instance";

export interface SigninReq {
  username?: string;
  email: string;
  password: string;
}

export interface SigninRes extends Omit<User, "username"> {
  username?: string;
}

const signup = async (body: SigninReq): Promise<SigninRes> => {
  const res = await instance.post(`/api/v1/auth/signup`, body);

  return res.data;
};

export default signup;
