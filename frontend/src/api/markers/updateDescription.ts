import instance from "../instance";

const updateDescription = async (desc: string, id: number) => {
  const formData = new FormData();

  formData.append("description", desc);

  try {
    const res = await instance.put(`/api/v1/markers/${id}`, formData, {
      withCredentials: true,
    });

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default updateDescription;
