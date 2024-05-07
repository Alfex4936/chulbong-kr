// // "use client";

// // import deleteFavorite from "@/api/favorite/deleteFavorite";
// // import setFavorite from "@/api/favorite/setFavorite";
// // import getMarker from "@/api/markers/getMarker";
// // import getWeather from "@/api/markers/getWeather";
// // import { useToast } from "@/components/ui/use-toast";
// // import { MOBILE_WIDTH } from "@/constants";
// // import useMapStore from "@/store/useMapStore";
// // import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
// // import useRoadviewStatusStore from "@/store/useRoadviewStatusStore";
// // import { Photo } from "@/types/Marker.types";
// // import { isAxiosError } from "axios";
// // import { useRouter } from "next/navigation";
// // import { useState } from "react";

// // const useSetMarker = () => {
// //   const router = useRouter();

// //   const { close: mobileMapClose } = useMobileMapOpenStore();
// //   const { toast } = useToast();
// //   const { map, clusterer } = useMapStore();
// //   const { open: openRoadview, setPosition: setRoadview } =
// //     useRoadviewStatusStore();

// //   const [bookmarkError, setBookmarkError] = useState(false);

// //   const set = async (lat: number, lng: number, markerId: number) => {
// //     console.log(1);
// //     if (!map || !clusterer) return;
// //     console.log(2);
// //     const imageSize = new window.kakao.maps.Size(39, 39);
// //     const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

// //     const activeMarkerImg = new window.kakao.maps.MarkerImage(
// //       "/activeMarker.svg",
// //       imageSize,
// //       imageOption
// //     );

// //     const skeletoncontent = document.createElement("div");
// //     skeletoncontent.className = "skeleton-overlay";

// //     const content = document.createElement("div");

// //     const skeletonOverlay = new window.kakao.maps.CustomOverlay({
// //       content: skeletoncontent,
// //       zIndex: 5,
// //     });

// //     const changeRoadviewlocation = async () => {
// //       setRoadview(lat, lng);
// //     };

// //     const copyTextToClipboard = async () => {
// //       const url = `${process.env.NEXT_PUBLIC_URL}/pullup/${markerId}`;
// //       try {
// //         await navigator.clipboard.writeText(url);
// //         toast({
// //           description: "링크 복사 완료",
// //         });
// //       } catch (err) {
// //         alert("잠시 후 다시 시도해 주세요!");
// //       }
// //     };

// //     console.log(markerId, lat, lng);

// //     const newMarker = new window.kakao.maps.Marker({
// //       map: map,
// //       position: new window.kakao.maps.LatLng(lat, lng),
// //       image: activeMarkerImg,
// //       title: markerId,
// //       zIndex: 4,
// //     });

// //     let markerLoading = false;
// //     let weatherLoading = false;

// //     window.kakao.maps.event.addListener(newMarker, "click", async () => {
// //       if (weatherLoading || markerLoading) return;
// //       content.innerHTML = "";
// //       const infoBox = /* HTML */ `
// //         <div id="overlay-top">
// //           <div id="overlay-weather">
// //             <div>
// //               <img id="overlay-weather-icon" />
// //             </div>
// //             <div id="overlay-weather-temp"></div>
// //           </div>
// //           <button id="overlay-close">닫기</button>
// //         </div>
// //         <div id="overlay-mid">
// //           <div id="overlay-info">
// //             <div id="overlay-title"></div>
// //             <div id="overlay-link">
// //               <button id="item-detail-link">상세보기</button>
// //               <button>정보 수정 제안</button>
// //             </div>
// //             <div class="empty-grow"></div>
// //             <div id="overlay-action">
// //               <button id="bookmark-button">
// //                 <div>
// //                   <img
// //                     id="bookmark-button-img"
// //                     src="/bookmark-02.svg"
// //                     alt="bookmark"
// //                   />
// //                 </div>
// //                 <div id="bookmark-text">북마크</div>
// //               </button>
// //               <button id="roadview-button">
// //                 <div>
// //                   <img src="/roadview.svg" alt="roadview" />
// //                 </div>
// //                 <div>거리뷰</div>
// //               </button>
// //               <button id="share-button">
// //                 <div>
// //                   <img src="/share-08.svg" alt="share" />
// //                 </div>
// //                 <div>공유</div>
// //               </button>
// //             </div>
// //           </div>
// //           <div id="overlay-image-container">
// //             <img id="overlay-image" />
// //           </div>
// //         </div>
// //       `;

// //       content.className = "overlay";
// //       content.innerHTML = infoBox;

// //       const overlay = new window.kakao.maps.CustomOverlay({
// //         content: content,
// //         zIndex: 5,
// //       });

// //       const latlng = new window.kakao.maps.LatLng(lat, lng);

// //       skeletonOverlay.setMap(map);
// //       skeletonOverlay.setPosition(latlng);

// //       // 마커 정보
// //       let description: string = "";
// //       let address: string = "";
// //       let favorited: boolean = false;
// //       let photos: Photo[] = [];
// //       let markerError = false;
// //       // 날씨 정보
// //       let iconImage: string = "";
// //       let temperature: string = "";
// //       let desc: string = "";
// //       let weatherError = false;
// //       // 북마크 정보
// //       let addBookmarkLoading = false;
// //       let addBookmarkError = false;
// //       let deleteBookmarkLoading = false;
// //       let deleteBookmarkError = false;

// //       const fetchMarker = async () => {
// //         markerLoading = true;
// //         try {
// //           const res = await getMarker(markerId);
// //           description = res.description;
// //           address = res.address as string;
// //           favorited = res.favorited as boolean;
// //           photos = res.photos as Photo[];
// //         } catch (error) {
// //           markerError = true;
// //           content.innerHTML = /* HTML */ `
// //             <div class="error-box">
// //               <span>잘못된 위치입니다. 잠시 후 다시 시도해 주세요.</span>
// //               <span><button id="error-close">닫기</button></span>
// //             </div>
// //           `;
// //           const errorCloseBtn = document.getElementById("error-close");
// //           errorCloseBtn?.addEventListener("click", () => {
// //             overlay.setMap(null);
// //           });
// //         } finally {
// //           markerLoading = false;
// //         }
// //       };

// //       const fetchWeather = async () => {
// //         weatherLoading = true;
// //         try {
// //           const res = await getWeather(lat, lng);
// //           iconImage = res.iconImage;
// //           temperature = res.temperature;
// //           desc = res.desc;
// //         } catch (error) {
// //           weatherError = true;
// //         } finally {
// //           weatherLoading = false;
// //         }
// //       };

// //       const addBookmark = async () => {
// //         addBookmarkLoading = true;
// //         try {
// //           const res = await setFavorite(markerId);
// //           return res;
// //         } catch (error) {
// //           if (isAxiosError(error)) {
// //             if (error.response?.status === 401) open();
// //           } else {
// //             toast({
// //               description: "잠시 후 다시 시도해 주세요",
// //             });
// //           }
// //           addBookmarkError = true;
// //           setBookmarkError(true);
// //         } finally {
// //           addBookmarkLoading = false;
// //         }
// //       };

// //       const deleteBookmark = async () => {
// //         deleteBookmarkLoading = true;
// //         try {
// //           const res = await deleteFavorite(markerId);
// //           return res;
// //         } catch (error) {
// //           deleteBookmarkError = true;
// //           setBookmarkError(true);
// //         } finally {
// //           deleteBookmarkLoading = false;
// //         }
// //       };

// //       await fetchMarker();
// //       await fetchWeather();

// //       skeletonOverlay.setMap(null);

// //       overlay.setMap(map);
// //       overlay.setPosition(latlng);

// //       // 오버레이 날씨 정보
// //       const weatherIconBox = document.getElementById(
// //         "overlay-weather-icon"
// //       ) as HTMLImageElement;
// //       if (weatherIconBox) {
// //         weatherIconBox.src = `${iconImage}` || "";
// //         weatherIconBox.alt = `${desc} || ""`;
// //       }

// //       const weatherTempBox = document.getElementById(
// //         "overlay-weather-temp"
// //       ) as HTMLDivElement;
// //       if (weatherTempBox) {
// //         weatherTempBox.innerHTML = `${temperature}℃`;
// //       }

// //       // 오버레이 주소 정보
// //       const addressBox = document.getElementById(
// //         "overlay-title"
// //       ) as HTMLDivElement;
// //       if (addressBox) {
// //         addressBox.innerHTML = description || "작성된 설명이 없습니다.";
// //       }

// //       // 오버레이 이미지 정보
// //       const imageContainer = document.getElementById(
// //         "overlay-image-container"
// //       ) as HTMLDivElement;
// //       if (imageContainer) {
// //         imageContainer.classList.add("on-loading");
// //       }
// //       const imageBox = document.getElementById(
// //         "overlay-image"
// //       ) as HTMLImageElement;
// //       if (imageBox) {
// //         imageBox.src = photos ? photos[0]?.photoUrl : "/metaimg.webp";
// //         imageBox.onload = () => {
// //           imageBox.style.display = "block";
// //           imageContainer.classList.remove("on-loading");
// //         };
// //       }

// //       // 오버레이 상세보기 링크
// //       const detailLink = document.getElementById(
// //         "item-detail-link"
// //       ) as HTMLAnchorElement;
// //       if (detailLink) {
// //         detailLink.style.cursor = "pointer";
// //         detailLink.addEventListener("click", () => {
// //           if (window.innerWidth <= MOBILE_WIDTH) {
// //             mobileMapClose();
// //           }
// //           router.push(`/pullup/${markerId}`);
// //         });
// //       }

// //       // 오버레이 북마크 버튼 이미지
// //       const bookmarkBtnImg = document.getElementById(
// //         "bookmark-button-img"
// //       ) as HTMLImageElement;
// //       if (bookmarkBtnImg) {
// //         bookmarkBtnImg.src = favorited
// //           ? "/bookmark-03.svg"
// //           : "/bookmark-02.svg";
// //       }

// //       // 오버레이 북마크 버튼 액션
// //       const bookmarkBtn = document.getElementById(
// //         "bookmark-button-img"
// //       ) as HTMLButtonElement;
// //       const bookmarkText = document.getElementById(
// //         "bookmark-text"
// //       ) as HTMLDivElement;
// //       if (bookmarkBtn && bookmarkText) {
// //         bookmarkBtn.addEventListener("click", async () => {
// //           if (addBookmarkLoading || deleteBookmarkLoading) return;
// //           bookmarkBtn.disabled = true;
// //           if (favorited) {
// //             bookmarkText.innerHTML = "취소중..";
// //             await deleteBookmark();
// //           } else if (!favorited) {
// //             bookmarkText.innerHTML = "저장중..";
// //             await addBookmark();
// //           }
// //           await fetchMarker();

// //           bookmarkText.innerHTML = "북마크";
// //           bookmarkBtnImg.src = favorited
// //             ? "/bookmark-03.svg"
// //             : "/bookmark-02.svg";

// //           bookmarkBtn.disabled = false;
// //         });
// //       }

// //       // 오보레이 로드뷰 버튼
// //       const roadviewButton = document.getElementById(
// //         "roadview-button"
// //       ) as HTMLButtonElement;
// //       if (roadviewButton) {
// //         roadviewButton.addEventListener("click", async () => {
// //           await changeRoadviewlocation();
// //           openRoadview();
// //         });
// //       }

// //       // 오버레이 공유 버튼
// //       const shareButton = document.getElementById(
// //         "share-button"
// //       ) as HTMLButtonElement;
// //       if (shareButton) {
// //         shareButton.addEventListener("click", copyTextToClipboard);
// //       }

// //       // 오버레이 닫기 이벤트 등록
// //       const closeBtnBox = document.getElementById(
// //         "overlay-close"
// //       ) as HTMLButtonElement;
// //       if (closeBtnBox) {
// //         closeBtnBox.onclick = () => {
// //           overlay.setMap(null);
// //         };
// //       }

// //       // 에러 오버레이 닫기
// //       const errorCloseBtn = document.getElementById("error-close");
// //       if (errorCloseBtn) {
// //         errorCloseBtn.onclick = () => {
// //           overlay.setMap(null);
// //         };
// //       }
// //     });

// //     console.log(newMarker);
// //     // newMarker.setMap(map);
// //     console.log(clusterer);
// //     clusterer.addMarker(newMarker);
// //     console.log(clusterer);
// //   };

// //   return { set };
// // };

// // export default useSetMarker;
// "use client";

// import deleteFavorite from "@/api/favorite/deleteFavorite";
// import setFavorite from "@/api/favorite/setFavorite";
// import getMarker from "@/api/markers/getMarker";
// import getWeather from "@/api/markers/getWeather";
// import { useToast } from "@/components/ui/use-toast";
// import { MOBILE_WIDTH } from "@/constants";
// import useMapStore from "@/store/useMapStore";
// import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
// import useRoadviewStatusStore from "@/store/useRoadviewStatusStore";
// import { Photo } from "@/types/Marker.types";
// import { isAxiosError } from "axios";
// import { useRouter } from "next/navigation";
// import { useState } from "react";

// const useSetMarker = () => {
//   const router = useRouter();

//   const { close: mobileMapClose } = useMobileMapOpenStore();
//   const { toast } = useToast();
//   const { map } = useMapStore();
//   const { open: openRoadview, setPosition: setRoadview } =
//     useRoadviewStatusStore();

//   const [bookmarkError, setBookmarkError] = useState(false);

//   const set = async (lat: number, lng: number, markerId: number) => {
//     console.log(1);
//     if (!map) return;
//     console.log(2);
//     const imageSize = new window.kakao.maps.Size(39, 39);
//     const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

//     const activeMarkerImg = new window.kakao.maps.MarkerImage(
//       "/activeMarker.svg",
//       imageSize,
//       imageOption
//     );

//     const skeletoncontent = document.createElement("div");
//     skeletoncontent.className = "skeleton-overlay";

//     const content = document.createElement("div");

//     const skeletonOverlay = new window.kakao.maps.CustomOverlay({
//       content: skeletoncontent,
//       zIndex: 5,
//     });

//     const changeRoadviewlocation = async () => {
//       setRoadview(lat, lng);
//     };

//     const copyTextToClipboard = async () => {
//       const url = `${process.env.NEXT_PUBLIC_URL}/pullup/${markerId}`;
//       try {
//         await navigator.clipboard.writeText(url);
//         toast({
//           description: "링크 복사 완료",
//         });
//       } catch (err) {
//         alert("잠시 후 다시 시도해 주세요!");
//       }
//     };

//     console.log(markerId, lat, lng);

//     const newMarker = new window.kakao.maps.Marker({
//       position: new window.kakao.maps.LatLng(lat, lng),
//       image: activeMarkerImg,
//       title: markerId,
//       zIndex: 4,
//     });

//     let markerLoading = false;
//     let weatherLoading = false;

//     window.kakao.maps.event.addListener(newMarker, "click", async () => {
//       if (weatherLoading || markerLoading) return;
//       content.innerHTML = "";
//       const infoBox = /* HTML */ `
//         <div id="overlay-top">
//           <div id="overlay-weather">
//             <div>
//               <img id="overlay-weather-icon" />
//             </div>
//             <div id="overlay-weather-temp"></div>
//           </div>
//           <button id="overlay-close">닫기</button>
//         </div>
//         <div id="overlay-mid">
//           <div id="overlay-info">
//             <div id="overlay-title"></div>
//             <div id="overlay-link">
//               <button id="item-detail-link">상세보기</button>
//               <button>정보 수정 제안</button>
//             </div>
//             <div class="empty-grow"></div>
//             <div id="overlay-action">
//               <button id="bookmark-button">
//                 <div>
//                   <img
//                     id="bookmark-button-img"
//                     src="/bookmark-02.svg"
//                     alt="bookmark"
//                   />
//                 </div>
//                 <div id="bookmark-text">북마크</div>
//               </button>
//               <button id="roadview-button">
//                 <div>
//                   <img src="/roadview.svg" alt="roadview" />
//                 </div>
//                 <div>거리뷰</div>
//               </button>
//               <button id="share-button">
//                 <div>
//                   <img src="/share-08.svg" alt="share" />
//                 </div>
//                 <div>공유</div>
//               </button>
//             </div>
//           </div>
//           <div id="overlay-image-container">
//             <img id="overlay-image" />
//           </div>
//         </div>
//       `;

//       content.className = "overlay";
//       content.innerHTML = infoBox;

//       const overlay = new window.kakao.maps.CustomOverlay({
//         content: content,
//         zIndex: 5,
//       });

//       const latlng = new window.kakao.maps.LatLng(lat, lng);

//       skeletonOverlay.setMap(map);
//       skeletonOverlay.setPosition(latlng);

//       // 마커 정보
//       let description: string = "";
//       let address: string = "";
//       let favorited: boolean = false;
//       let photos: Photo[] = [];
//       let markerError = false;
//       // 날씨 정보
//       let iconImage: string = "";
//       let temperature: string = "";
//       let desc: string = "";
//       let weatherError = false;
//       // 북마크 정보
//       let addBookmarkLoading = false;
//       let addBookmarkError = false;
//       let deleteBookmarkLoading = false;
//       let deleteBookmarkError = false;

//       const fetchMarker = async () => {
//         markerLoading = true;
//         try {
//           const res = await getMarker(markerId);
//           description = res.description;
//           address = res.address as string;
//           favorited = res.favorited as boolean;
//           photos = res.photos as Photo[];
//         } catch (error) {
//           markerError = true;
//           content.innerHTML = /* HTML */ `
//             <div class="error-box">
//               <span>잘못된 위치입니다. 잠시 후 다시 시도해 주세요.</span>
//               <span><button id="error-close">닫기</button></span>
//             </div>
//           `;
//           const errorCloseBtn = document.getElementById("error-close");
//           errorCloseBtn?.addEventListener("click", () => {
//             overlay.setMap(null);
//           });
//         } finally {
//           markerLoading = false;
//         }
//       };

//       const fetchWeather = async () => {
//         weatherLoading = true;
//         try {
//           const res = await getWeather(lat, lng);
//           iconImage = res.iconImage;
//           temperature = res.temperature;
//           desc = res.desc;
//         } catch (error) {
//           weatherError = true;
//         } finally {
//           weatherLoading = false;
//         }
//       };

//       const addBookmark = async () => {
//         addBookmarkLoading = true;
//         try {
//           const res = await setFavorite(markerId);
//           return res;
//         } catch (error) {
//           if (isAxiosError(error)) {
//             if (error.response?.status === 401) open();
//           } else {
//             toast({
//               description: "잠시 후 다시 시도해 주세요",
//             });
//           }
//           addBookmarkError = true;
//           setBookmarkError(true);
//         } finally {
//           addBookmarkLoading = false;
//         }
//       };

//       const deleteBookmark = async () => {
//         deleteBookmarkLoading = true;
//         try {
//           const res = await deleteFavorite(markerId);
//           return res;
//         } catch (error) {
//           deleteBookmarkError = true;
//           setBookmarkError(true);
//         } finally {
//           deleteBookmarkLoading = false;
//         }
//       };

//       await fetchMarker();
//       await fetchWeather();

//       skeletonOverlay.setMap(null);

//       overlay.setMap(map);
//       overlay.setPosition(latlng);

//       // 오버레이 날씨 정보
//       const weatherIconBox = document.getElementById(
//         "overlay-weather-icon"
//       ) as HTMLImageElement;
//       if (weatherIconBox) {
//         weatherIconBox.src = `${iconImage}` || "";
//         weatherIconBox.alt = `${desc} || ""`;
//       }

//       const weatherTempBox = document.getElementById(
//         "overlay-weather-temp"
//       ) as HTMLDivElement;
//       if (weatherTempBox) {
//         weatherTempBox.innerHTML = `${temperature}℃`;
//       }

//       // 오버레이 주소 정보
//       const addressBox = document.getElementById(
//         "overlay-title"
//       ) as HTMLDivElement;
//       if (addressBox) {
//         addressBox.innerHTML = description || "작성된 설명이 없습니다.";
//       }

//       // 오버레이 이미지 정보
//       const imageContainer = document.getElementById(
//         "overlay-image-container"
//       ) as HTMLDivElement;
//       if (imageContainer) {
//         imageContainer.classList.add("on-loading");
//       }
//       const imageBox = document.getElementById(
//         "overlay-image"
//       ) as HTMLImageElement;
//       if (imageBox) {
//         imageBox.src = photos ? photos[0]?.photoUrl : "/metaimg.webp";
//         imageBox.onload = () => {
//           imageBox.style.display = "block";
//           imageContainer.classList.remove("on-loading");
//         };
//       }

//       // 오버레이 상세보기 링크
//       const detailLink = document.getElementById(
//         "item-detail-link"
//       ) as HTMLAnchorElement;
//       if (detailLink) {
//         detailLink.style.cursor = "pointer";
//         detailLink.addEventListener("click", () => {
//           if (window.innerWidth <= MOBILE_WIDTH) {
//             mobileMapClose();
//           }
//           router.push(`/pullup/${markerId}`);
//         });
//       }

//       // 오버레이 북마크 버튼 이미지
//       const bookmarkBtnImg = document.getElementById(
//         "bookmark-button-img"
//       ) as HTMLImageElement;
//       if (bookmarkBtnImg) {
//         bookmarkBtnImg.src = favorited
//           ? "/bookmark-03.svg"
//           : "/bookmark-02.svg";
//       }

//       // 오버레이 북마크 버튼 액션
//       const bookmarkBtn = document.getElementById(
//         "bookmark-button-img"
//       ) as HTMLButtonElement;
//       const bookmarkText = document.getElementById(
//         "bookmark-text"
//       ) as HTMLDivElement;
//       if (bookmarkBtn && bookmarkText) {
//         bookmarkBtn.addEventListener("click", async () => {
//           if (addBookmarkLoading || deleteBookmarkLoading) return;
//           bookmarkBtn.disabled = true;
//           if (favorited) {
//             bookmarkText.innerHTML = "취소중..";
//             await deleteBookmark();
//           } else if (!favorited) {
//             bookmarkText.innerHTML = "저장중..";
//             await addBookmark();
//           }
//           await fetchMarker();

//           bookmarkText.innerHTML = "북마크";
//           bookmarkBtnImg.src = favorited
//             ? "/bookmark-03.svg"
//             : "/bookmark-02.svg";

//           bookmarkBtn.disabled = false;
//         });
//       }

//       // 오보레이 로드뷰 버튼
//       const roadviewButton = document.getElementById(
//         "roadview-button"
//       ) as HTMLButtonElement;
//       if (roadviewButton) {
//         roadviewButton.addEventListener("click", async () => {
//           await changeRoadviewlocation();
//           openRoadview();
//         });
//       }

//       // 오버레이 공유 버튼
//       const shareButton = document.getElementById(
//         "share-button"
//       ) as HTMLButtonElement;
//       if (shareButton) {
//         shareButton.addEventListener("click", copyTextToClipboard);
//       }

//       // 오버레이 닫기 이벤트 등록
//       const closeBtnBox = document.getElementById(
//         "overlay-close"
//       ) as HTMLButtonElement;
//       if (closeBtnBox) {
//         closeBtnBox.onclick = () => {
//           overlay.setMap(null);
//         };
//       }

//       // 에러 오버레이 닫기
//       const errorCloseBtn = document.getElementById("error-close");
//       if (errorCloseBtn) {
//         errorCloseBtn.onclick = () => {
//           overlay.setMap(null);
//         };
//       }
//     });

//     return newMarker;
//   };

//   return { set };
// };

// export default useSetMarker;
