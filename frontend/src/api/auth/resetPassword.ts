import axios from "axios";

const resetPassword = async (token: string, password: string) => {
  const formData = new FormData();

  formData.append("password", password);
  formData.append("token", token);

  try {
    const res = await axios.post(`/api/v1/auth/reset-password`, formData);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default resetPassword;
