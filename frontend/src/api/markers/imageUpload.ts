import axios from "axios";

const imageUpload = async (file: File) => {
  const formData = new FormData();
  const token = JSON.parse(localStorage.getItem("user") as string).state.user
    .token;

  formData.append("file", file);

  try {
    const res = await axios.post(
      `${import.meta.env.VITE_LOCAL_URL}/api/v1/markers/upload`,
      formData,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    return res;
  } catch (error) {
    if (axios.isAxiosError(error) && error.response) {
      console.error(
        `이미지 업로드 실패: ${error.response.status} - ${error.response.data.error}`
      );

      return error.response.data.error;
    } else {
      console.error(`이미지 업로드 실패: ${error}`);
    }
  }
};

export default imageUpload;
