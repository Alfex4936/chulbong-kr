import axios from "axios";

const deleteMarker = async (id: number) => {
  const token = JSON.parse(localStorage.getItem("user") as string).state.user
    .token;

  try {
    const res = await axios.delete(`/api/v1/markers/${id}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    return res;
  } catch (error) {
    throw error;
  }
};

export default deleteMarker;
