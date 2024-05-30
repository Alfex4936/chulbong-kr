// import { type MarkerRes } from "@/api/markers/getAllMarker";

// interface WeatherRes {
//   temperature: string;
//   desc: string;
//   humidity: string;
//   rainfall: string;
//   snowfall: string;
//   iconImage: string;
// }

// const overlayContent = () => {
//   const content = document.createElement("div");
//   content.className = "overlay";
//   return content;
// };

// const overlayInfo = (
//   weather: WeatherRes,
//   address: string,
//   markerId: number
// ) => {
//   const { iconImage, temperature, desc } = weather;

//   const infoBoxasd = `
//     <div>
//       <div id="overlay-close">닫기</div>
//       <div>날씨</div>
//       <div>{address}</div>
//       <div>
//         <span>상세보기</span>
//         <span>정보 수정 제안</span>
//       </div>
//       <div>
//         <span>저장</span>
//         <span>거리뷰</span>
//         <span>공유</span>
//       </div>
//     </div>`;

//   // weahter container
//   const infoBox = document.createElement("div");
//   infoBox.className = "overlayInfo";

//   // weather icon container
//   const weatherBox = document.createElement("div");

//   const weatherIconWrap = document.createElement("div");
//   const weatherIcon = document.createElement("img");
//   weatherIcon.src = iconImage;
//   weatherIcon.alt = desc;
//   weatherIconWrap.appendChild(weatherIcon);

//   //  weather temperature container
//   const weatherTemp = document.createElement("div");
//   weatherTemp.innerHTML = `${temperature}℃`;

//   // append weather info
//   weatherBox.appendChild(weatherBox);
//   weatherBox.appendChild(weatherTemp);

//   infoBox.appendChild(weatherBox);

//   // address
//   const addressBox = document.createElement("h1");
//   addressBox.innerHTML = address;

//   // info link
//   const linkBox = document.createElement("div");
//   const detailLink = document.createElement("a");
//   detailLink.href = `/chullbong/${markerId}`;
//   const suggestLink = document.createElement("a");
//   suggestLink.href = "/suggestion";

//   linkBox.appendChild(detailLink);
//   linkBox.appendChild(suggestLink);

//   infoBox.appendChild(linkBox);

//   // action button
//   const bookmarkBox = document.createElement("div");
//   const bookmarkIcon = document.createElement("img");
//   bookmarkIcon.src = "/bookmark.svg";

//   const roadviewBox = document.createElement("div");
//   const shareBox = document.createElement("div");

//   return infoBox;
// };

// const overlayImage = (imgUrl: string) => {
//   const imageBox = document.createElement("div");
//   const image = document.createElement("img");
//   image.src = imgUrl || "/metaimg.webp";

//   return imageBox;
// };

// const generateOverlay = (
//   imgUrl: string,
//   data: MarkerRes,
//   weather: WeatherRes
// ) => {
//   const { markerId, address } = data;

//   const content = overlayContent();
//   const info = overlayInfo(
//     weather,
//     address || "주소 정보가 없습니다.",
//     markerId
//   );
//   const image = overlayImage(imgUrl);

//   content.appendChild(info);
//   content.appendChild(image);

//   return content;
// };

// export default generateOverlay;
