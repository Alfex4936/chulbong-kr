import instance from "../instance";

const getMyReports = async () => {
  const res = await instance.get(`/api/v1/users/reports`);

  return res.data;
};

export default getMyReports;
