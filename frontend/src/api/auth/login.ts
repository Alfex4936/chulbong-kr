import type { User } from "@/types/User.types";
import axios from "axios";

export interface LoginReq {
  email: string;
  password: string;
}

export interface LoginRes {
  token: string;
  user: User;
}

const login = async (body: LoginReq): Promise<LoginRes> => {
  try {
    const res = await axios.post(`/api/v1/auth/login`, body);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default login;
