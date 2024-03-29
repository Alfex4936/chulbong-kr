import instance from "../instance";

const logout = async () => {
  try {
    const res = await instance.post(`/api/v1/auth/logout`, {
      withCredentials: true,
    });

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default logout;
