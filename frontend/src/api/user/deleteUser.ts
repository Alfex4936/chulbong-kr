import instance from "../instance";

const deleteUser = async () => {
  try {
    const res = await instance.delete(`/api/v1/users/me`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default deleteUser;
