import instance from "../instance";

const checkAdmin = async () => {
  const res = await instance.get(`/api/v1/admin`);

  return res.data;
};

export default checkAdmin;
