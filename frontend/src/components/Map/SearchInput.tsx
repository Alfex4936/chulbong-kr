import LocationOnIcon from "@mui/icons-material/LocationOn";
import SearchIcon from "@mui/icons-material/Search";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { useEffect, useState } from "react";
import getSearchLoation from "../../api/kakao/getSearchLoation";
import useInput from "../../hooks/useInput";
import * as Styled from "./SearchInput.style";
import type { KakaoMap } from "../../types/KakaoMap.types";
import ExpandLessIcon from "@mui/icons-material/ExpandLess";

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
}

const SearchInput = ({ map }: Props) => {
  const searchInput = useInput("");
  const [places, setPlaces] = useState<KakaoPlace[] | null>(null);
  const [isResult, setIsResult] = useState(false);

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
    map.setCenter(moveLatLon);
    map?.setLevel(2);
  };

  return (
    <div
      style={{
        position: "absolute",
        top: "10px",
        left: "10px",
        zIndex: "10",
        width: "50%",
      }}
    >
      <Styled.InputWrap>
        <Styled.SearchInput
          type="text"
          name="search"
          placeholder="ex) 햄버거 맛집, 수원, 잠실역, 남산 타워"
          value={searchInput.value}
          onChange={searchInput.onChange}
        />
        <Tooltip title="검색" arrow disableInteractive>
          <IconButton aria-label="send" onClick={handleSearch}>
            <SearchIcon />
          </IconButton>
        </Tooltip>
      </Styled.InputWrap>
      {isResult && (
        <Styled.Result>
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
          <Tooltip title="닫기" arrow disableInteractive>
            <IconButton
              onClick={() => {
                setIsResult(false);
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
