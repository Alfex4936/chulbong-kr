import axios from "axios";

const updateDescription = async (desc: string, id: number) => {
  const formData = new FormData();

  formData.append("description", desc);

  try {
    const res = await axios.put(`/api/v1/markers/${id}`, formData, {
      withCredentials: true,
    });

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default updateDescription;
