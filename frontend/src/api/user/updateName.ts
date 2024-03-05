import axios from "axios";

// interface MyInfo {
//   userId: number;
//   username: string;
//   email: string;
// }

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
