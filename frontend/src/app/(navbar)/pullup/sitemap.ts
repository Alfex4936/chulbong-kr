import getAllMarker from "@/api/markers/getAllMarker";
import { BASE_URL } from "@/constants";
import { MetadataRoute } from "next";

const sitemap = async (): Promise<MetadataRoute.Sitemap> => {
  const products = await getAllMarker();
  return products.map((marker) => ({
    url: `${BASE_URL}/pullup/${marker.markerId}`,
  }));
};

export default sitemap;
