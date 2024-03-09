import axios from "axios";

const adminCheck = async () => {
  try {
    const res = await axios.get(`/api/v1/admin`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default adminCheck;
