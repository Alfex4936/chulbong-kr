import axios from "axios";

const getSearchLoation = async (query: string) => {
  try {
    const res = await axios.get(
      `https://dapi.kakao.com/v2/local/search/keyword.json?query=${query}&page=1&size=5`,
      {
        headers: {
          Authorization: `KakaoAK ${import.meta.env.VITE_KAKAO_API_KEY}`,
        },
      }
    );

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default getSearchLoation;
