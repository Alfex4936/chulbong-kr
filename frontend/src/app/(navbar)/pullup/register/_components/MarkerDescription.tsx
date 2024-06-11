"use client";

import getMarker from "@/api/markers/getMarker";
import ErrorMessage from "@/components/atom/ErrorMessage";
import LoadingSpinner from "@/components/atom/LoadingSpinner";
import MarkerOverlay from "@/components/layout/MarkerOverlay";
import RQProvider from "@/components/provider/RQProvider";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { MAP_LAT_DIF } from "@/constants";
import useInput from "@/hooks/common/useInput";
import useMapControl from "@/hooks/common/useMapControl";
import useSetFacilities from "@/hooks/mutation/marker/useSetFacilities";
import useUploadMarker from "@/hooks/mutation/marker/useUploadMarker";
import useReportMarker from "@/hooks/mutation/report/useReportMarker";
import useMapStatusStore from "@/store/useMapStatusStore";
import useMapStore from "@/store/useMapStore";
import useUploadFormDataStore from "@/store/useUploadFormDataStore";
import { isAxiosError } from "axios";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { createRoot } from "react-dom/client";

interface Props {
  desc?: string;
  markerId?: number;
  isReport?: boolean;
}

const MarkerDescription = ({ desc, markerId, isReport = false }: Props) => {
  const descriptionValue = useInput(desc ? desc : "");

  const router = useRouter();

  const { mutateAsync: setFacilities } = useSetFacilities();
  const { mutateAsync: uploadMarker } = useUploadMarker();
  const { clusterer, map, setMarkers, markers, setOverlay } = useMapStore();
  const { setPosition } = useMapStatusStore();
  const { moveLocation } = useMapControl();

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const [errorMessage, setErrorMessage] = useState("");

  const { imageForm, latitude, longitude, resetData, facilities } =
    useUploadFormDataStore();

  const {
    mutate: report,
    error: reportError,
    isError,
  } = useReportMarker(markerId as number);

  useEffect(() => {
    resetData();
  }, []);

  const handleSubmit = async () => {
    if (isReport && markerId) {
      setLoading(true);
      try {
        const marker = await getMarker(markerId);
        let data;
        if (latitude !== 0 && longitude !== 0) {
          data = {
            markerId: markerId,
            description: descriptionValue.value,
            photos: imageForm.map((image) => image.file) as File[],
            latitude: marker.latitude,
            longitude: marker.longitude,
            newLatitude: latitude,
            newLongitude: longitude,
          };
        } else {
          data = {
            markerId: markerId,
            description: descriptionValue.value,
            photos: imageForm.map((image) => image.file) as File[],
            latitude: marker.latitude,
            longitude: marker.longitude,
          };
        }

        if (imageForm.length <= 0) {
          setErrorMessage("이미지를 1개 이상 등록해 주세요");
          setLoading(false);
          return;
        }

        report(data);
      } catch (error) {
        setErrorMessage("잠시 후 다시 시도해 주세요.");
        setLoading(false);
      }

      return;
    }

    const data = {
      description: descriptionValue.value,
      photos: imageForm.map((image) => image.file) as File[],
      latitude: latitude,
      longitude: longitude,
    };

    setLoading(true);

    try {
      const result = await uploadMarker(data);
      await setFacilities({
        markerId: result.markerId,
        facilities: [
          {
            facilityId: 1,
            quantity: facilities.철봉,
          },
          {
            facilityId: 2,
            quantity: facilities.평행봉,
          },
        ],
      });

      const imageSize = new window.kakao.maps.Size(39, 39);
      const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

      const activeMarkerImg = new window.kakao.maps.MarkerImage(
        "/selectedMarker.svg",
        imageSize,
        imageOption
      );

      const overlayDiv = document.createElement("div");
      overlayDiv.classList.add("overlay_1");
      const root = createRoot(overlayDiv);

      const newMarker = new window.kakao.maps.Marker({
        position: new window.kakao.maps.LatLng(latitude, longitude),
        image: activeMarkerImg,
        title: result.markerId,
        zIndex: 4,
      });

      window.kakao.maps.event.addListener(newMarker, "click", async () => {
        map?.setLevel(3);

        const moveLatLon = new window.kakao.maps.LatLng(
          (latitude as number) + MAP_LAT_DIF,
          longitude
        );

        map?.panTo(moveLatLon);
        setPosition((latitude as number) + MAP_LAT_DIF, longitude);

        if (document.getElementsByClassName("overlay_1")[0]) {
          document.getElementsByClassName("overlay_1")[0].remove();
        }

        const latlng = new window.kakao.maps.LatLng(latitude, longitude);
        const overlay = new window.kakao.maps.CustomOverlay({
          position: latlng,
          content: overlayDiv,
          zIndex: 11,
        });

        setOverlay(overlay);

        overlay.setMap(map);

        const closeOverlay = () => {
          overlay.setMap(null);
        };

        const goDetail = () => {
          router.push(`/pullup/${result.markerId}`);
        };
        const goReport = () => {
          router.push(`/pullup/${result.markerId}/report`);
        };

        root.render(
          <RQProvider>
            <MarkerOverlay
              markerId={result.markerId}
              closeOverlay={closeOverlay}
              goDetail={goDetail}
              goReport={goReport}
              lat={result.latitude}
              lng={result.longitude}
            />
          </RQProvider>
        );
      });

      if (!markers) return;
      const newMarkers = [...markers, newMarker];
      clusterer?.addMarker(newMarker);

      setMarkers(newMarkers);
      moveLocation({ lat: latitude, lng: longitude, isfilter: true });

      router.push(`/pullup/${result.markerId}`);
    } catch (error) {
      if (isAxiosError(error)) {
        if (error.response?.status === 401) {
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
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isAxiosError(reportError)) {
      if (reportError.response?.status === 406) {
        setErrorMessage("기존 위치와 너무 많은 거리가 차이납니다.");
      }

      setLoading(false);
    }
  }, [isError]);

  return (
    <div>
      <div className="flex flex-col mb-5">
        <Input
          className="text-base"
          type="text"
          id="description"
          placeholder={desc ? "" : "설명 입력"}
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
          {loading ? (
            <LoadingSpinner size="xs" />
          ) : markerId && isReport ? (
            "제안 요청"
          ) : (
            "등록하기"
          )}
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
      <ErrorMessage text={errorMessage} />
    </div>
  );
};

export default MarkerDescription;
