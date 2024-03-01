import axios from "axios";

interface Props {
  email: string;
  code: string;
}

const verifyCode = async ({ email, code }: Props): Promise<string> => {
  const formData = new FormData();

  formData.append("email", email);
  formData.append("token", code);

  try {
    const res = await axios.post(`/api/v1/auth/verify-email/confirm`, formData);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default verifyCode;
