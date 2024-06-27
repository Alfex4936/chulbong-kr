const getRegion = (): {
  getCode: (name: string) => string;
  getTitle: (code: string) => string;
} => {
  const getCode = (name: string) => {
    if (name === "제주특별자치도") return "jj";
    else if (name === "전남") return "jn";
    else if (name === "전북특별자치도") return "jb";
    else if (name === "경남") return "gn";
    else if (name === "경북") return "gb";
    else if (name === "대구") return "dg";
    else if (name === "울산") return "us";
    else if (name === "충북") return "cb";
    else if (name === "충남") return "cn";
    else if (name === "대전") return "dj";
    else if (name === "강원특별자치도") return "gw";
    else if (name === "경기") return "gg";
    else if (name === "서울") return "so";
    else if (name === "인천") return "ic";
    else if (name === "부산") return "bs";
    else return "";
  };

  const getTitle = (code: string) => {
    if (code === "jj") return "제주도 채팅방";
    else if (code === "jn") return "전라남도 채팅방";
    else if (code === "jb") return "전북특별자치도 채팅방";
    else if (code === "gn") return "경상남도 채팅방";
    else if (code === "gb") return "경상북도 채팅방";
    else if (code === "dg") return "대구 채팅방";
    else if (code === "us") return "울산 채팅방";
    else if (code === "cb") return "충청북도 채팅방";
    else if (code === "cn") return "충청남code도 채팅방";
    else if (code === "dj") return "대전 채팅방";
    else if (code === "gw") return "강원도 채팅방";
    else if (code === "gg") return "경기도 채팅방";
    else if (code === "so") return "서울 채팅방";
    else if (code === "ic") return "인천 채팅방";
    else if (code === "bs") return "부산 채팅방";
    else return "";
  };

  return { getCode, getTitle };
};

export default getRegion;
