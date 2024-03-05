import axios from "axios";

const requestPasswordReset = async (email: string) => {
  const formData = new FormData();

  formData.append("email", email);

  try {
    const res = await axios.post(
      `/api/v1/auth/request-password-reset`,
      formData
    );

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default requestPasswordReset;
