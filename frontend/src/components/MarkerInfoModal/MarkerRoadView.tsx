import Button from "@mui/material/Button";
import { useEffect, useRef } from "react";
import * as Styled from "./MarkerRoadView.style";

interface Props {
  lat: number;
  lng: number;
  setIsRoadView: React.Dispatch<React.SetStateAction<boolean>>;
  setIsRoadViewError: React.Dispatch<React.SetStateAction<boolean>>;
}

const MarkerRoadView = ({
  lat,
  lng,
  setIsRoadView,
  setIsRoadViewError,
}: Props) => {
  const roadViewRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const roadview = new window.kakao.maps.Roadview(roadViewRef.current);
    const roadviewClient = new window.kakao.maps.RoadviewClient();

    const position = new window.kakao.maps.LatLng(lat, lng);

    roadviewClient.getNearestPanoId(position, 50, (panoId: number) => {
      if (panoId === null) {
        setIsRoadViewError(true);
        setIsRoadView(false);
        return;
      }

      roadview.setPanoId(panoId, position);
    });
  }, []);

  return (
    <Styled.Container>
      <Styled.Exit>
        <Button
          sx={{
            color: "#fff",
            fontSize: "1rem",
            "&:hover": {
              backgroundColor: "rgba(0,0,0,0.5)",
            },
          }}
          onClick={() => {
            setIsRoadView(false);
          }}
        >
          닫기
        </Button>
      </Styled.Exit>
      <Styled.RoadViewContainer ref={roadViewRef} />
    </Styled.Container>
  );
};

export default MarkerRoadView;
