import instance from "@/api/instance";
import { type MarkerRes } from "@/api/markers/getAllMarker";
import { BASE_URL } from "@/constants";
import { MetadataRoute } from "next";

const getAllMarker = async (): Promise<MarkerRes[]> => {
  const res = await instance.get(`${process.env.NEXT_PUBLIC_BASE_URL}/markers`);

  return res.data;
};

const sitemap = async (): Promise<MetadataRoute.Sitemap> => {
  const pullup = await getAllMarker();
  const pullupMap = pullup.map((marker) => ({
    url: `${BASE_URL}/pullup/${marker.markerId}`,
  }));

  const routesMap = ["", "/home", "/mypage", "/search", "/signin", "/signup"].map(
    (route) => ({
      url: `${BASE_URL}${route}`,
    })
  );

  return [...routesMap, ...pullupMap];
};

export default sitemap;
