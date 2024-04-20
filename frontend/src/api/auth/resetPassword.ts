import instance from "../instance";

interface Props {
  token: string;
  password: string;
}

const resetPassword = async ({ token, password }: Props) => {
  const formData = new FormData();

  formData.append("password", password);
  formData.append("token", token);

  const res = await instance.post(`/api/v1/auth/reset-password`, formData);

  return res.data;
};

export default resetPassword;
