import LocationOnIcon from "@mui/icons-material/LocationOn";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { ComponentProps, forwardRef } from "react";
import activeMarkerImage from "../../assets/images/cb1.webp";
import selectedMarkerImage from "../../assets/images/cb3.webp";
import useGetFavorites from "../../hooks/query/favorites/useGetFavorites";
import useCurrentMarkerStore from "../../store/useCurrentMarkerStore";
import useMapPositionStore from "../../store/useMapPositionStore";
import type { KakaoMap, KakaoMarker } from "../../types/KakaoMap.types";
import * as Styled from "./FavoriteMarker.style";

interface Props extends ComponentProps<"div"> {
  markers: KakaoMarker[];
  map: KakaoMap;
}

const FavoriteMarker = forwardRef(({ map, markers, ...props }: Props, ref) => {
  const mapPosition = useMapPositionStore();
  const currentMarkerState = useCurrentMarkerStore();

  const { data, isLoading, isError } = useGetFavorites();

  if (data?.length === 0) {
    return <div style={{ padding: "1rem" }}>등록한 위치가 없습니다.</div>;
  }

  const handleMove = (lat: number, lon: number) => {
    const moveLatLon = new window.kakao.maps.LatLng(lat, lon);

    mapPosition.setPosition(lat, lon);
    mapPosition.setLevel(1);

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

  if (isLoading) {
    return <Styled.ListSkeleton />;
  }

  if (isError)
    return <div style={{ padding: "1rem" }}>등록한 위치가 없습니다.</div>;

  return (
    <Styled.Container ref={ref as React.RefObject<HTMLDivElement>} {...props}>
      {data?.map((marker) => (
        <Styled.ListContainer key={marker.markerId}>
          <Styled.MarkerList>
            <Styled.MarkerListTop>
              <Styled.Description>{marker.description}</Styled.Description>
              <Styled.AddressText>{marker.addr}</Styled.AddressText>
            </Styled.MarkerListTop>
            <div>
              <Tooltip title="이동" arrow disableInteractive>
                <IconButton
                  onClick={() => {
                    handleMove(marker.latitude, marker.longitude);
                    filtering(marker.markerId);
                    currentMarkerState.setMarker(marker.markerId);
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
