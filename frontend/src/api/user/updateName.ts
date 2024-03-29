import instance from "../instance";

const updateName = async (name: string) => {
  try {
    const res = await instance.patch(`/api/v1/users/me`, {
      username: name,
    });

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default updateName;
