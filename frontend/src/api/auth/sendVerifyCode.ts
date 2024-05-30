import instance from "../instance";

const sendVerifyCode = async (email: string): Promise<string> => {
  const formData = new FormData();

  formData.append("email", email);

  const res = await instance.post(`/api/v1/auth/verify-email/send`, formData);

  return res.data;
};

export default sendVerifyCode;
