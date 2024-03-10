import LocationOnIcon from "@mui/icons-material/LocationOn";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { ComponentProps, forwardRef } from "react";
import useGetFavorites from "../../hooks/query/favorites/useGetFavorites";
import useMapPositionStore from "../../store/useMapPositionStore";
import type { KakaoMap } from "../../types/KakaoMap.types";
import * as Styled from "./FavoriteMarker.style";

interface Props extends ComponentProps<"div"> {
  map: KakaoMap;
}

const FavoriteMarker = forwardRef(({ map, ...props }: Props, ref) => {
  const mapPosition = useMapPositionStore();

  const { data, isLoading, isError } = useGetFavorites();

  if (data?.length === 0) {
    return <div style={{ padding: "1rem" }}>등록한 위치가 없습니다.</div>;
  }

  if (isLoading) {
    return <Styled.ListSkeleton />;
  }

  if (isError)
    return <div style={{ padding: "1rem" }}>등록한 위치가 없습니다.</div>;

  const handleMove = (lat: number, lon: number) => {
    const moveLatLon = new window.kakao.maps.LatLng(lat, lon);

    mapPosition.setPosition(lat, lon);
    mapPosition.setLevel(1);

    map.setCenter(moveLatLon);
    map?.setLevel(1);
  };

  return (
    <Styled.Container ref={ref as React.RefObject<HTMLDivElement>} {...props}>
      {data?.map((marker) => (
        <Styled.ListContainer key={marker.markerId}>
          <Styled.MarkerList>
            <Styled.MarkerListTop>
              <p style={{ flexGrow: "1", textAlign: "left" }}>
                {marker.description || "설명 없음"}
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
        </Styled.ListContainer>
      ))}
      <Styled.LoadList />
    </Styled.Container>
  );
});

export default FavoriteMarker;
