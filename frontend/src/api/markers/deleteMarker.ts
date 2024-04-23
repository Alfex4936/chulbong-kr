import instance from "../instance";

const deleteMarker = async (id: number) => {
  const res = await instance.delete(`/api/v1/markers/${id}`, {
    withCredentials: true,
  });

  return res;
};

export default deleteMarker;
