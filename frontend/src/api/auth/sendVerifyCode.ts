import instance from "../instance";

const sendVerifyCode = async (email: string): Promise<string> => {
  const formData = new FormData();

  formData.append("email", email);

  try {
    const res = await instance.post(`/api/v1/auth/verify-email/send`, formData);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default sendVerifyCode;
