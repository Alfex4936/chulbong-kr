import ExpandLessIcon from "@mui/icons-material/ExpandLess";
import LocationOnIcon from "@mui/icons-material/LocationOn";
import SearchIcon from "@mui/icons-material/Search";
import Button from "@mui/material/Button";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { RefObject, useEffect, useRef, useState } from "react";
import getSearchLoation from "../../api/kakao/getSearchLoation";
import useInput from "../../hooks/useInput";
import useMapPositionStore from "../../store/useMapPositionStore";
import useOnBoardingStore from "../../store/useOnBoardingStore";
import type { KakaoMap, KakaoMarker } from "../../types/KakaoMap.types";
import AroundMarker from "../AroundMarker/AroundMarker";
import Ranking from "../Ranking/Ranking";
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
  markerInfoModal: boolean;
  setIsAround: React.Dispatch<React.SetStateAction<boolean>>;
}

const SearchInput = ({
  map,
  markers,
  aroundMarkerRef,
  isAround,
  markerInfoModal,
  setIsAround,
}: Props) => {
  const mapPosition = useMapPositionStore();
  const onBoardingState = useOnBoardingStore();

  const searchInput = useInput("");

  const searchInputRef = useRef<HTMLInputElement>(null);

  const [places, setPlaces] = useState<KakaoPlace[] | null>(null);
  const [isResult, setIsResult] = useState(false);

  const [curTab, setCurTab] = useState<number>(0);

  useEffect(() => {
    const handleKeyDownClose = (event: KeyboardEvent) => {
      if (markerInfoModal) {
        searchInputRef.current?.blur();
        return;
      }
      if (event.key === "/") {
        event.preventDefault();
        searchInputRef.current?.focus();
      } else if (event.key === "Escape") {
        searchInput.reset();
        searchInputRef.current?.blur();
      }
    };

    window.addEventListener("keydown", handleKeyDownClose);

    return () => {
      window.removeEventListener("keydown", handleKeyDownClose);
    };
  }, [markerInfoModal]);

  useEffect(() => {
    if (!onBoardingState.isOnBoarding) {
      setIsResult(false);
      return;
    }

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
    } else if (onBoardingState.step === 12 || onBoardingState.step === 13) {
      setCurTab(1);
    } else {
      setIsResult(false);
      setCurTab(0);
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

  const tabs = [
    {
      title: "주변 검색",
      content: (
        <AroundMarker map={map} ref={aroundMarkerRef} markers={markers} />
      ),
    },
    {
      title: "랭킹",
      content: <Ranking map={map} />,
    },
  ];

  return (
    <div
      style={{
        flexGrow: "1",
        zIndex: onBoardingState.step === 8 ? 1005 : 200,
      }}
    >
      <Styled.InputWrap>
        <Styled.SearchInput
          ref={searchInputRef}
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
            <div>
              <Styled.TabContainer>
                {tabs.map((tab, index) => {
                  return (
                    <Button
                      key={index}
                      sx={{
                        width: "50%",
                        color: index === curTab ? "#6b73db" : "#333",
                      }}
                      onClick={() => {
                        setCurTab(index);
                      }}
                    >
                      {tab.title}
                    </Button>
                  );
                })}
              </Styled.TabContainer>
              {curTab !== null && <div>{tabs[curTab].content}</div>}
            </div>
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
                setCurTab(0);
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
