import instance from "../instance";

const saveOffline = async (lat: number, lng: number): Promise<Blob> => {
  const res = await instance.get(
    `/api/v1/markers/save-offline?latitude=${lat}&longitude=${lng}`
  );

  return res.data;
};

export default saveOffline;
