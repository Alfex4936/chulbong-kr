"use client";

import useMyinfoData from "@/hooks/query/user/useMyinfoData";
import Unauthenticated from "./_component/Unauthenticated";
import useLogout from "@/hooks/mutation/auth/useLogout";

const MypageClient = () => {
  const { data: myInfo, isError } = useMyinfoData();
  const { mutate: logout } = useLogout();

  // console.log(myInfo);

  if (!myInfo || isError) return <Unauthenticated />;
  return (
    <div>
      <div>{myInfo?.username}</div>
      <div>{myInfo?.email}</div>

      <button
        onClick={() => {
          logout();
        }}
      >
        로그아웃
      </button>
    </div>
  );
};

export default MypageClient;
