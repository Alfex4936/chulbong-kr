import axios from "axios";

const updateName = async (name: string) => {
  try {
    const res = await axios.patch(`/api/v1/users/me`, {
      username: name,
    });

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default updateName;
