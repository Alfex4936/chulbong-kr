import LocationOnIcon from "@mui/icons-material/LocationOn";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { useState } from "react";
import useMarkerAriaRankingData from "../../hooks/query/marker/useMarkerAriaRankingData";
import useMarkerRankingData from "../../hooks/query/marker/useMarkerRankingData";
import useMapPositionStore from "../../store/useMapPositionStore";
import type { KakaoMap } from "../../types/KakaoMap.types";
import ActionButton from "../ActionButton/ActionButton";
import * as Styled from "./Ranking.style";
import type { RankingInfo } from "../../api/markers/markerRanking";

import { useTranslation } from 'react-i18next';
import '../../i18n';

const TOP10_RANKING_NUMBER = 0;
const AROUND_RANKING_NUMBER = 1;

interface Props {
  map: KakaoMap;
}

const Ranking = ({ map }: Props) => {
  const { t } = useTranslation();

  const mapPosition = useMapPositionStore();

  const {
    data: topRanking,
    refetch: fetchTotopRankingRanking,
    isLoading: topRankingLoading,
    isFetching: topRankingFetching,
  } = useMarkerRankingData();
  const {
    data: topRankingAria,
    refetch: fetchTopRankingAria,
    isLoading: topRankingAriaLoading,
    isFetching: topRankingAriaFetching,
  } = useMarkerAriaRankingData(mapPosition.lat, mapPosition.lng);

  const [curRanking, setCurRanking] = useState(TOP10_RANKING_NUMBER);

  const hanldeClickTopRanking = () => {
    fetchTotopRankingRanking();
    setCurRanking(TOP10_RANKING_NUMBER);
  };

  const hanldeClickTopRankingAria = () => {
    fetchTopRankingAria();
    setCurRanking(AROUND_RANKING_NUMBER);
  };

  const handleMove = (lat: number, lon: number) => {
    const moveLatLon = new window.kakao.maps.LatLng(lat, lon);

    mapPosition.setPosition(lat, lon);
    mapPosition.setLevel(2);

    map.setCenter(moveLatLon);
    map?.setLevel(2);
  };
  return (
    <Styled.Container>
      <Styled.MessageRed>
        {t("ranking.description")}
      </Styled.MessageRed>
      <Styled.ButtonContainer>
        <ActionButton
          bg={curRanking === TOP10_RANKING_NUMBER ? "black" : "gray"}
          onClick={hanldeClickTopRanking}
        >
          TOP 10
        </ActionButton>
        <ActionButton
          bg={curRanking === AROUND_RANKING_NUMBER ? "black" : "gray"}
          onClick={hanldeClickTopRankingAria}
        >
          {t("ranking.title")}
        </ActionButton>
      </Styled.ButtonContainer>
      {topRankingLoading ||
        topRankingAriaLoading ||
        topRankingFetching ||
        topRankingAriaFetching ? (
        <Styled.ListSkeleton />
      ) : (
        <>
          {curRanking === TOP10_RANKING_NUMBER && (
            <>
              {topRanking && topRanking?.length > 0 ? (
                <>
                  {topRanking?.map((item: RankingInfo, index) => {
                    return (
                      <Styled.ResultItem key={item.makerId}>
                        <span>{index + 1}{t("ranking.position")}</span>
                        <span>{item.address}</span>
                        <span>
                          <Tooltip title={t("map.move")} arrow disableInteractive>
                            <IconButton
                              onClick={() => {
                                handleMove(item.latitude, item.longitude);
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
                        </span>
                      </Styled.ResultItem>
                    );
                  })}
                </>
              ) : (
                <div>{t("ranking.noLocations")}</div>
              )}
            </>
          )}

          {curRanking === AROUND_RANKING_NUMBER && (
            <>
              {topRankingAria && topRankingAria?.length > 0 ? (
                <>
                  {topRankingAria?.map((item: RankingInfo, index) => {
                    return (
                      <Styled.ResultItem key={item.makerId}>
                        <span>{index + 1} {t("ranking.position")}</span>
                        <span>{item.address}</span>
                        <span>
                          <Tooltip title={t("map.move")} arrow disableInteractive>
                            <IconButton
                              onClick={() => {
                                handleMove(item.latitude, item.longitude);
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
                        </span>
                      </Styled.ResultItem>
                    );
                  })}
                </>
              ) : (
                <div>{t("ranking.noLocations")}</div>
              )}
            </>
          )}
        </>
      )}
    </Styled.Container>
  );
};

export default Ranking;
