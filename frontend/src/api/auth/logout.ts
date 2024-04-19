import instance from "../instance";

const logout = async () => {
  const res = await instance.post(`/api/v1/auth/logout`, {
    withCredentials: true,
  });

  return res.data;
};

export default logout;
