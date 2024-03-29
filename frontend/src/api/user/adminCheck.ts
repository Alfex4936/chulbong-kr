import instance from "../instance";

const adminCheck = async () => {
  try {
    const res = await instance.get(`/api/v1/admin`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default adminCheck;
