import TravelExploreIcon from "@mui/icons-material/TravelExplore";
import Button from "@mui/material/Button";
import Tooltip from "@mui/material/Tooltip";
import { useEffect, useRef, useState } from "react";
import useOnBoardingStore from "../../store/useOnBoardingStore";
import type { KakaoMap, KakaoMarker } from "../../types/KakaoMap.types";
import * as Styled from "./MapHeader.style";
import SearchInput from "./SearchInput";

interface Props {
  markers: KakaoMarker[];
  map: KakaoMap | null;
}

const MapHeader = ({ markers, map }: Props) => {
  const onBoardingState = useOnBoardingStore();

  const [isAround, setIsAround] = useState(false);
  const aroundMarkerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!onBoardingState.isOnBoarding) {
      setIsAround(false);
      return;
    }

    if (onBoardingState.step >= 11 && onBoardingState.step <= 13) {
      setIsAround(true);
    } else {
      setIsAround(false);
    }
  }, [onBoardingState.step]);

  const handleArroundMarkerScroll = () => {
    if (aroundMarkerRef.current) {
      aroundMarkerRef.current.scrollTop = 0;
    }
  };

  return (
    <Styled.Container>
      <SearchInput
        map={map as KakaoMap}
        markers={markers}
        aroundMarkerRef={aroundMarkerRef}
        isAround={isAround}
        setIsAround={setIsAround}
      />

      <Tooltip
        title={isAround ? "스크롤 위로" : "주변 검색 / 랭킹"}
        arrow
        disableInteractive
      >
        <Button
          sx={{
            boxShadow:
              "rgba(50, 50, 93, 0.25) 0px 2px 5px -1px, rgba(0, 0, 0, 0.3) 0px 1px 3px -1px",

            backgroundColor: "#fff",
            color: "#000",

            zIndex: onBoardingState.step === 10 ? "10000" : "90",

            borderRadius: ".5rem",

            "&:hover": {
              backgroundColor: "#888",
              color: "#fff",
            },
          }}
          onClick={() => {
            if (onBoardingState.step === 10) return;
            setIsAround(true);
            handleArroundMarkerScroll();
          }}
        >
          <TravelExploreIcon />
        </Button>
      </Tooltip>
    </Styled.Container>
  );
};

export default MapHeader;
