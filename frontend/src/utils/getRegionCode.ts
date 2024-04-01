const getRegion = (
  name?: string,
  code?: string
): { getCode: () => number; getTitle: () => string } => {
  const getCode = () => {
    if (name === "제주특별자치도") return 1;
    else if (name === "전남") return 2;
    else if (name === "전북특별자치도") return 3;
    else if (name === "경남") return 4;
    else if (name === "경북") return 5;
    else if (name === "대구") return 6;
    else if (name === "울산") return 7;
    else if (name === "충북") return 8;
    else if (name === "충남") return 9;
    else if (name === "대전") return 10;
    else if (name === "강원특별자치도") return 11;
    else if (name === "경기") return 12;
    else if (name === "서울") return 13;
    else if (name === "인천") return 14;
    else if (name === "부산") return 15;
    else return 30;
  };

  const getTitle = () => {
    if (code === "1") return "제주도 채팅방";
    else if (name === "2") return "전라남도 채팅방";
    else if (name === "3") return "전북특별자치도 채팅방";
    else if (name === "4") return "경상남도 채팅방";
    else if (name === "5") return "경상북도 채팅방";
    else if (name === "6") return "대구 채팅방";
    else if (name === "7") return "울산 채팅방";
    else if (name === "8") return "충청북도 채팅방";
    else if (name === "9") return "충청남도 채팅방";
    else if (name === "10") return "대전 채팅방";
    else if (name === "11") return "강원도 채팅방";
    else if (name === "12") return "경기도 채팅방";
    else if (name === "13") return "서울 채팅방";
    else if (name === "14") return "인천 채팅방";
    else if (name === "15") return "부산 채팅방";
    else return "";
  };

  return { getCode, getTitle };
};

export default getRegion;
