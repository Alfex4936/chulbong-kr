"use client";

import logout from "@/api/auth/logout";
import LoadingSpinner from "@/components/atom/LoadingSpinner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import useInput from "@/hooks/common/useInput";
import useUploadFormDataStore from "@/store/useUploadFormDataStore";
import { isAxiosError } from "axios";
import { useEffect, useState } from "react";

const MarkerDescription = () => {
  const descriptionValue = useInput("");
  const [error, setError] = useState("");

  const [loading, setLoading] = useState(false);

  const { imageForm, latitude, longitude, resetData } =
    useUploadFormDataStore();

  useEffect(() => {
    resetData();
  }, []);

  const handleSubmit = async () => {
    const data = {
      description: descriptionValue.value,
      photos: imageForm.map((image) => image.file) as File[],
      latitude: latitude,
      longitude: longitude,
    };

    setLoading(true);

    try {
      console.log(data);
      //   TODO: 업로드 로직 작성하기
      //   const result = await uploadMarker(data);
      //   await setFacilities({
      //     markerId: result.markerId,
      //     facilities: [
      //       {
      //         facilityId: 1,
      //         quantity: chulbong,
      //       },
      //       {
      //         facilityId: 2,
      //         quantity: penghang,
      //       },
      //     ],
      //   });
      //   await filtering();
      //   const imageSize = new window.kakao.maps.Size(39, 39);
      //   const imageOption = { offset: new window.kakao.maps.Point(27, 45) };
      //   const selectedMarkerImg = new window.kakao.maps.MarkerImage(
      //     selectedMarkerImage,
      //     imageSize,
      //     imageOption
      //   );
      //   const newMarker = new window.kakao.maps.Marker({
      //     map: map,
      //     position: new window.kakao.maps.LatLng(
      //       formState.latitude,
      //       formState.longitude
      //     ),
      //     image: selectedMarkerImg,
      //     title: result.markerId,
      //     zIndex: 4,
      //   });
      //   window.kakao.maps.event.addListener(newMarker, "click", () => {
      //     setMarkerInfoModal(true);
      //     setCurrentMarkerInfo({
      //       markerId: result.markerId,
      //     } as MarkerInfo);
      //   });
      //   setMarkers((prev) => {
      //     const copy = [...prev];
      //     copy.push(newMarker);
      //     return copy;
      //   });
      //   setState(false);
      //   setIsMarked(false);
      //   clusterer.addMarker(newMarker);
      //   marker?.setMap(null);
      //   currentMarkerState.setMarker(result.markerId);
    } catch (error) {
      if (isAxiosError(error)) {
        if (error.response?.status === 401) {
          await logout();
          setError("인증이 만료 되었습니다. 다시 로그인 해주세요!");
        } else if (error.response?.status === 409) {
          setError("주변에 이미 철봉이 있습니다!");
        } else if (error.response?.status === 403) {
          setError("대한민국에서만 등록 가능합니다!");
        } else if (error.response?.status === 400) {
          setError("입력을 확인해 주세요!");
        } else {
          setError("잠시 후 다시 시도해 주세요!");
        }
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <div className="flex flex-col mb-5">
        <Input
          type="text"
          id="description"
          placeholder="설명 입력"
          maxLength={70}
          value={descriptionValue.value}
          onChange={(e) => {
            if (descriptionValue.value.length >= 70) {
              setError("70자 이내로 작성해 주세요!");
            } else {
              setError("");
            }
            descriptionValue.handleChange(e);
          }}
        />
        <div
          data-testid="file-error"
          className="mt-1 text-center text-sm text-red"
        >
          {error}
        </div>
      </div>

      <div>
        <Button
          className="border-grey border bg-transparent dark:text-grey hover:bg-white-tp-light hover:border-transparent"
          onClick={handleSubmit}
          disabled={loading}
        >
          {loading ? <LoadingSpinner size="sm" /> : "등록하기"}
        </Button>
        {loading && imageForm.length > 0 ? (
          <div
            data-testid="file-error"
            className="mt-1 text-center text-sm text-red"
          >
            이미지를 등록하는 중입니다. 잠시만 기다려 주세요!
          </div>
        ) : null}
      </div>
    </div>
  );
};

export default MarkerDescription;
