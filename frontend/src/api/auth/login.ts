import type { User } from "@/types/User.types";
import instance from "../instance";

export interface LoginReq {
  email: string;
  password: string;
}

export interface LoginRes {
  token: string;
  user: User;
}

const login = async (body: LoginReq): Promise<LoginRes> => {
  const res = await instance.post(`/api/v1/auth/login`, body);

  return res.data;
};

export default login;
