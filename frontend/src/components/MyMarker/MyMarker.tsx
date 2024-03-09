import LocationOnIcon from "@mui/icons-material/LocationOn";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { useEffect, useRef } from "react";
import useGetMyMarker from "../../hooks/query/useGetMyMarker";
import type { KakaoMap } from "../../types/KakaoMap.types";
import * as Styled from "./MyMarker.style";

interface Props {
  map: KakaoMap;
}

const MyMarker = ({ map }: Props) => {
  const { data, fetchNextPage, hasNextPage, isLoading, isError, isFetching } =
    useGetMyMarker();

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

  if (data?.pages.length === 0) {
    return <div>등록한 위치가 없습니다.</div>;
  }

  if (isLoading) {
    return <Styled.ListSkeleton />;
  }
  if (isError)
    return <div style={{ padding: "1rem" }}>등록한 위치가 없습니다.</div>;

  const handleMove = (lat: number, lon: number) => {
    const moveLatLon = new window.kakao.maps.LatLng(lat, lon);
    map.setCenter(moveLatLon);
    map?.setLevel(1);
  };

  return (
    <Styled.Container>
      {data?.pages.map((page, i) => (
        <Styled.ListContainer key={i}>
          {page.markers.map((marker) => (
            <Styled.MarkerList key={marker.markerId}>
              <Styled.MarkerListTop>
                <p style={{ flexGrow: "1", textAlign: "left" }}>
                  {marker.description}
                </p>
                <Styled.AddressText>{marker.addr}</Styled.AddressText>
              </Styled.MarkerListTop>
              <div>
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
              </div>
            </Styled.MarkerList>
          ))}
        </Styled.ListContainer>
      ))}
      <Styled.LoadList />
      {hasNextPage && <Styled.ListSkeleton ref={boxRef} />}
    </Styled.Container>
  );
};

export default MyMarker;
