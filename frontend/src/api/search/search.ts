import instance from "../instance";

interface Markers {
  address: string;
  markerId: number;
}

export interface SearchRes {
  took: number;
  markers: Markers[];
}

const search = async (query: string): Promise<SearchRes> => {
  const res = await instance.get(`/api/v1/search/marker?term=${query}`);

  return res.data;
};

export default search;
