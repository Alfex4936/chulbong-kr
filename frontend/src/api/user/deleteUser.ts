import instance from "../instance";

const deleteUser = async () => {
  const res = await instance.delete(`/api/v1/users/me`);

  return res.data;
};

export default deleteUser;
