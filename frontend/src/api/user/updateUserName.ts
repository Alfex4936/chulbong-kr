import instance from "../instance";

const updateUserName = async (name: string) => {
  const res = await instance.patch(`/api/v1/users/me`, {
    username: name,
  });

  return res.data;
};

export default updateUserName;
