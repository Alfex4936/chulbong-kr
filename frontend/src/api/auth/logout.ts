import axios from "axios";

const logout = async () => {
  try {
    const res = await axios.post(`/api/v1/auth/logout`, {
      withCredentials: true,
    });

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default logout;
