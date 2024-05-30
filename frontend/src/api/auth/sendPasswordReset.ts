import instance from "../instance";

const sendPasswordReset = async (email: string) => {
  const formData = new FormData();

  formData.append("email", email);

  const res = await instance.post(
    `/api/v1/auth/request-password-reset`,
    formData
  );

  return res.data;
};

export default sendPasswordReset;
