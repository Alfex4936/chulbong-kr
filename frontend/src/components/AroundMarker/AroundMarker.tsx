import LocationOnIcon from "@mui/icons-material/LocationOn";
import SearchIcon from "@mui/icons-material/Search";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { useEffect, useRef, useState } from "react";
import useGetCloseMarker from "../../hooks/query/useGetCloseMarker";
import useMapPositionStore from "../../store/useMapPositionStore";
import type { KakaoMap } from "../../types/KakaoMap.types";
import * as Styled from "./AroundMarker.style";

interface Props {
  map: KakaoMap;
}

const AroundMarker = ({ map }: Props) => {
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
    refetch();
  };

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setDistance(Number(event.target.value));
  };

  const handleMove = (lat: number, lon: number) => {
    const moveLatLon = new window.kakao.maps.LatLng(lat, lon);
    map.setCenter(moveLatLon);
    map?.setLevel(1);
  };

  return (
    <Styled.Container>
      <Styled.RangeContainer>
        <p style={{ fontSize: ".8rem" }}>주변 {distance}m</p>
        <input
          type="range"
          min="100"
          max="1000"
          step="100"
          value={distance}
          onChange={handleChange}
        />
        <Tooltip title="검색" arrow disableInteractive>
          <IconButton
            onClick={handleSearch}
            aria-label="delete"
            sx={{
              color: "#333",
              width: "30px",
              height: "30px",
            }}
          >
            <SearchIcon sx={{ fontSize: 22 }} />
          </IconButton>
        </Tooltip>
      </Styled.RangeContainer>
      {isLoading && (
        <Styled.ListSkeleton>
          <div />
          <div style={{ flexGrow: "1" }} />
          <div />
        </Styled.ListSkeleton>
      )}
      {data?.pages.map((page, i) => (
        <Styled.ListContainer key={i}>
          {page.markers?.map((marker) => (
            <Styled.MarkerList key={marker.markerId}>
              <p style={{ flexGrow: "1", textAlign: "left" }}>
                {marker.description}
              </p>
              <Tooltip title="이동" arrow disableInteractive>
                <IconButton
                  onClick={() => {
                    handleMove(marker.latitude, marker.longitude);
                  }}
                  aria-label="delete"
                  sx={{
                    color: "#333",
                    width: "25px",
                    height: "25px",
                  }}
                >
                  <LocationOnIcon sx={{ fontSize: 18 }} />
                </IconButton>
              </Tooltip>
            </Styled.MarkerList>
          ))}
        </Styled.ListContainer>
      ))}
      <Styled.LoadList />
      {data?.pages[0].markers === null && (
        <div style={{ padding: "1rem" }}>주변에 철봉이 없습니다.</div>
      )}
      {isError && (
        <div style={{ padding: "1rem" }}>잠시 후 다시 시도해 주세요</div>
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
};

export default AroundMarker;
