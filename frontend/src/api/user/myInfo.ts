import instance from "../instance";

export interface MyInfo {
  userId: number;
  username: string;
  email: string;
}

const myInfo = async (): Promise<MyInfo> => {
  const res = await instance.get(`/api/v1/users/me`);

  return res.data;
};

export default myInfo;
