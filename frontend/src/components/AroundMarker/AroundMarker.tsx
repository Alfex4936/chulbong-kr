import LocationOnIcon from "@mui/icons-material/LocationOn";
import SearchIcon from "@mui/icons-material/Search";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { useQueryClient } from "@tanstack/react-query";
import { ComponentProps, forwardRef, useEffect, useRef, useState } from "react";
import activeMarkerImage from "../../assets/images/cb1.webp";
import selectedMarkerImage from "../../assets/images/cb3.webp";
import useGetCloseMarker from "../../hooks/query/useGetCloseMarker";
import useCurrentMarkerStore from "../../store/useCurrentMarkerStore";
import useMapPositionStore from "../../store/useMapPositionStore";
import type { KakaoMap, KakaoMarker } from "../../types/KakaoMap.types";
import * as Styled from "./AroundMarker.style";

import { useTranslation } from 'react-i18next';
import '../../i18n';

interface Props extends ComponentProps<"div"> {
  markers: KakaoMarker[];
  map: KakaoMap;
}

const AroundMarker = forwardRef(({ markers, map, ...props }: Props, ref) => {
  const currentMarkerState = useCurrentMarkerStore();

  const queryClient = useQueryClient();

  const positionState = useMapPositionStore();
  const [distance, setDistance] = useState(100);

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isLoading,
    isError,
    isFetching,
    refetch,
  } = useGetCloseMarker({
    lat: positionState.lat,
    lon: positionState.lng,
    distance: distance,
  });
  const { t } = useTranslation();

  const boxRef = useRef(null);

  useEffect(() => {
    const currentRef = boxRef.current;

    const observer = new IntersectionObserver(
      (entries) => {
        const entry = entries[0];
        if (entry.isIntersecting) {
          if (!isFetching && hasNextPage) {
            fetchNextPage();
          }
        }
      },
      { threshold: 0.8 }
    );

    if (currentRef) {
      observer.observe(currentRef);
    }

    return () => {
      if (currentRef) {
        observer.unobserve(currentRef);
      }
    };
  }, [isFetching, hasNextPage, fetchNextPage]);

  const handleSearch = () => {
    queryClient.resetQueries({ queryKey: ["closeMarker", distance] });
    refetch();
  };

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setDistance(Number(event.target.value));
  };

  const handleMove = (lat: number, lon: number) => {
    const moveLatLon = new window.kakao.maps.LatLng(lat, lon);

    positionState.setPosition(lat, lon);
    positionState.setLevel(1);

    map.setCenter(moveLatLon);
    map?.setLevel(1);
  };

  const filtering = async (markerId: number) => {
    const imageSize = new window.kakao.maps.Size(39, 39);
    const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

    const selectedMarkerImg = new window.kakao.maps.MarkerImage(
      selectedMarkerImage,
      imageSize,
      imageOption
    );
    const activeMarkerImg = new window.kakao.maps.MarkerImage(
      activeMarkerImage,
      imageSize,
      imageOption
    );

    const marker = markers.find((value) => Number(value.Gb) === markerId);

    markers.forEach((marker) => {
      marker?.setImage(activeMarkerImg);
    });

    marker?.setImage(selectedMarkerImg);
  };

  return (
    <Styled.Container ref={ref as React.RefObject<HTMLDivElement>} {...props}>
      <div>
        <Styled.MessageRed>
          {t("around.searchInfo")}
        </Styled.MessageRed>
        <Styled.RangeContainer>
          <p>{t("around.within")} {distance}m</p>
          <div style={{ flexGrow: "1" }}>
            <input
              type="range"
              min="100"
              max="5000"
              step="100"
              value={distance}
              onChange={handleChange}
            />
          </div>
          <Tooltip title={t("map.search")} arrow disableInteractive>
            <IconButton
              onClick={handleSearch}
              aria-label="delete"
              sx={{
                color: "#333",
                width: "30px",
                height: "30px",
              }}
              disabled={isFetching}
            >
              <SearchIcon sx={{ fontSize: 22 }} />
            </IconButton>
          </Tooltip>
        </Styled.RangeContainer>
      </div>
      {isLoading ? (
        <Styled.ListSkeleton />
      ) : (
        <>
          {data?.pages.map((page, i) => (
            <Styled.ListContainer key={i}>
              {page.markers?.map((marker) => {
                return (
                  <Styled.MarkerList key={marker.markerId}>
                    <Styled.MarkerListTop>
                      <Styled.DescriptionWrap>
                        <Styled.Distance>
                          ({~~marker.distance}m)
                        </Styled.Distance>
                        <Styled.Description>
                          {marker.description || t("around.noDesc")}
                        </Styled.Description>
                      </Styled.DescriptionWrap>
                      <Styled.AddressText>
                        {marker.address ||
                          t("around.loadingAddr")}
                      </Styled.AddressText>
                    </Styled.MarkerListTop>
                    <div>
                      <Tooltip title={t("map.move")} arrow disableInteractive>
                        <IconButton
                          onClick={() => {
                            handleMove(marker.latitude, marker.longitude);
                            filtering(marker.markerId);
                            currentMarkerState.setMarker(marker.markerId);
                          }}
                          aria-label="move"
                          sx={{
                            color: "#333",
                            width: "25px",
                            height: "25px",
                          }}
                        >
                          <LocationOnIcon sx={{ fontSize: 18 }} />
                        </IconButton>
                      </Tooltip>
                    </div>
                  </Styled.MarkerList>
                );
              })}
            </Styled.ListContainer>
          ))}
        </>
      )}

      <Styled.LoadList />
      {data?.pages[0].markers === null && (
        <div style={{ padding: "1rem" }}>{t("around.noPullup")}</div>
      )}
      {isError && (
        <div style={{ padding: "1rem" }}>{t("around.tryAgain")}</div>
      )}
      {hasNextPage && (
        <Styled.ListSkeleton ref={boxRef}>
          <div />
          <div style={{ flexGrow: "1" }} />
          <div />
        </Styled.ListSkeleton>
      )}
    </Styled.Container>
  );
});

export default AroundMarker;
