import axios from "axios";

const getSearchLoation = async (query: string) => {
  const res = await axios.get(
    `https://dapi.kakao.com/v2/local/search/keyword.json?query=${query}&page=1&size=5`,
    {
      headers: {
        Authorization: `KakaoAK ${process.env.NEXT_PUBLIC_KAK}`,
      },
    }
  );

  return res.data;
};

export default getSearchLoation;
