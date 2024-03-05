import axios from "axios";

const deleteUser = async () => {
  try {
    const res = await axios.delete(`/api/v1/users/me`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default deleteUser;
