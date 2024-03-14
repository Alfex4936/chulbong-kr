import ExpandLessIcon from "@mui/icons-material/ExpandLess";
import LocationOnIcon from "@mui/icons-material/LocationOn";
import SearchIcon from "@mui/icons-material/Search";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { RefObject, useEffect, useState } from "react";
import getSearchLoation from "../../api/kakao/getSearchLoation";
import useInput from "../../hooks/useInput";
import useMapPositionStore from "../../store/useMapPositionStore";
import useOnBoardingStore from "../../store/useOnBoardingStore";
import type { KakaoMap, KakaoMarker } from "../../types/KakaoMap.types";
import AroundMarker from "../AroundMarker/AroundMarker";
import * as Styled from "./SearchInput.style";

interface KakaoPlace {
  address_name: string;
  category_group_code: string;
  category_group_name: string;
  category_name: string;
  distance: string;
  id: string;
  phone: string;
  place_name: string;
  place_url: string;
  road_address_name: string;
  x: string;
  y: string;
}

interface Props {
  map: KakaoMap;
  markers: KakaoMarker[];
  aroundMarkerRef: RefObject<HTMLDivElement>;
  isAround: boolean;
  setIsAround: React.Dispatch<React.SetStateAction<boolean>>;
}

const SearchInput = ({
  map,
  markers,
  aroundMarkerRef,
  isAround,
  setIsAround,
}: Props) => {
  const mapPosition = useMapPositionStore();
  const onBoardingState = useOnBoardingStore();

  const searchInput = useInput("");
  const [places, setPlaces] = useState<KakaoPlace[] | null>(null);
  const [isResult, setIsResult] = useState(false);

  useEffect(() => {
    if (!onBoardingState.isOnBoarding) return;

    const fetch = async () => {
      try {
        const result = await getSearchLoation("남산 타워");
        setPlaces(result.documents);
      } catch (error) {
        console.error(error);
      }
    };

    if (onBoardingState.step === 9) {
      setIsResult(true);
      fetch();
    } else {
      setIsResult(false);
    }
  }, [onBoardingState.step]);

  useEffect(() => {
    if (searchInput.value === "") {
      setIsResult(false);
      return;
    }

    setIsResult(true);
    const fetch = async () => {
      try {
        const result = await getSearchLoation(searchInput.value);
        setPlaces(result.documents);
      } catch (error) {
        console.error(error);
      }
    };

    fetch();
  }, [searchInput.value]);

  const handleSearch = () => {
    if (searchInput.value === "") {
      setIsResult(false);
      return;
    }
    setIsResult(true);

    const fetch = async () => {
      try {
        const result = await getSearchLoation(searchInput.value);
        setPlaces(result.documents);
      } catch (error) {
        console.log(error);
      }
    };

    fetch();
  };

  const handleMove = (lat: number, lon: number) => {
    const moveLatLon = new window.kakao.maps.LatLng(lat, lon);

    mapPosition.setPosition(lat, lon);
    mapPosition.setLevel(2);

    map.setCenter(moveLatLon);
    map?.setLevel(2);
  };

  return (
    <div
      style={{
        flexGrow: "1",
        zIndex: onBoardingState.step === 8 ? 1005 : 200,
      }}
    >
      <Styled.InputWrap>
        <Styled.SearchInput
          type="text"
          name="search"
          disabled={onBoardingState.step === 8}
          placeholder="ex) 햄버거 맛집, 수원, 잠실역, 남산 타워"
          value={
            onBoardingState.step === 8 || onBoardingState.step === 9
              ? "남산 타워"
              : searchInput.value
          }
          onChange={(e) => {
            searchInput.onChange(e);
            setIsAround(false);
          }}
        />
        <Tooltip title="검색" arrow disableInteractive>
          <IconButton aria-label="send" onClick={handleSearch}>
            <SearchIcon />
          </IconButton>
        </Tooltip>
      </Styled.InputWrap>
      {(isResult || isAround) && (
        <Styled.Result>
          {isAround ? (
            <AroundMarker map={map} ref={aroundMarkerRef} markers={markers} />
          ) : (
            <>
              {places?.map((place) => {
                return (
                  <Styled.ResultItem key={place.id}>
                    <div>
                      <span>{place.place_name}</span>
                      <span>({place.address_name})</span>
                    </div>
                    <Tooltip title="이동" arrow disableInteractive>
                      <IconButton
                        onClick={() => {
                          // console.log(Number(place.y), Number(place.x));
                          handleMove(Number(place.y), Number(place.x));
                        }}
                        aria-label="move"
                        sx={{
                          color: "#444",
                          width: "25px",
                          height: "25px",
                        }}
                      >
                        <LocationOnIcon sx={{ fontSize: 18 }} />
                      </IconButton>
                    </Tooltip>
                  </Styled.ResultItem>
                );
              })}
            </>
          )}
          <Tooltip title="닫기" arrow disableInteractive>
            <IconButton
              onClick={() => {
                setIsResult(false);
                setIsAround(false);
                searchInput.reset();
              }}
              aria-label="move"
              sx={{
                color: "#444",
                width: "25px",
                height: "25px",
              }}
            >
              <ExpandLessIcon />
            </IconButton>
          </Tooltip>
        </Styled.Result>
      )}
    </div>
  );
};

export default SearchInput;
