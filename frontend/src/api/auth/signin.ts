import axios from "axios";
import { User } from "@/types/User.types";

export interface SigninReq {
  username?: string;
  email: string;
  password: string;
}

export interface SigninRes extends Omit<User, "username"> {
  username?: string;
}

const signin = async (body: SigninReq): Promise<SigninRes> => {
  try {
    const res = await axios.post(`/api/v1/auth/signup`, body);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default signin;
