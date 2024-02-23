import axios from "axios";

const DeleteMarker = async (id: number) => {
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
    if (axios.isAxiosError(error) && error.response) {
      console.error(
        `삭제 실패: ${error.response.status} - ${error.response.data.error}`
      );

      return error.response.data.error;
    } else {
      console.error(`삭제 실패: ${error}`);
    }
  }
};

export default DeleteMarker;
