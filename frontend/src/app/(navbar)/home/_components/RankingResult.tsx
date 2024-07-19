"use client";

import { type RankingInfo } from "@/api/markers/getMarkerRanking";
import EmojiHoverButton from "@/components/atom/EmojiHoverButton";
import MarkerListItem from "@/components/atom/MarkerListItem";
import { Skeleton } from "@/components/ui/skeleton";
import useMyGps from "@/hooks/common/useMyGps";
import useAreaMarkerRankingData from "@/hooks/query/marker/useAreaMarkerRankingData";
import useMarkerRankingData from "@/hooks/query/marker/useMarkerRankingData";
import useMapStatusStore from "@/store/useMapStatusStore";
import useMapStore from "@/store/useMapStore";
import useScrollButtonStore from "@/store/useScrollButtonStore";
import useTabStore from "@/store/useTabStore";
import getAddress, { type AddressInfo } from "@/utils/getAddress";
import { useEffect, useState } from "react";

const CHUNK_SIZE = 10;

const RankingResult = () => {
  const { centerMapOnCurrentPositionAsync } = useMyGps();
  const { setIsActive } = useScrollButtonStore();

  const { map } = useMapStore();
  const { lat, lng } = useMapStatusStore();
  const { curTab, setDisable } = useTabStore();

  const { data: top10 } = useMarkerRankingData();
  const {
    data: areaRanking,
    isError: areaError,
    isLoading: areaLoading,
    refetch: areaRankingRefetch,
    isFetching: isAreaFetching,
  } = useAreaMarkerRankingData(lat, lng, curTab === "주변");

  const [address, setAddress] = useState("");
  const [sliceData, setSliceData] = useState<RankingInfo[][] | null>(null);

  const [pageNumber, setPageNumber] = useState(0);
  const [viewData, setViewData] = useState<RankingInfo[]>([]);

  useEffect(() => {
    if (!top10 || top10.length === 0) {
      setDisable(1);
      return;
    }

    const sliceArray = () => {
      const result = [];
      for (let i = 0; i < top10.length; i += CHUNK_SIZE) {
        const chunk = top10.slice(i, i + CHUNK_SIZE);
        result.push(chunk);
      }
      return result;
    };

    const chunked = sliceArray();
    setSliceData(chunked);
    setViewData(chunked[0]);
  }, [top10]);

  useEffect(() => {
    if (!sliceData || !viewData.length) return;

    if (pageNumber === 0) {
      setViewData(sliceData[0]);
      setIsActive(false);
    } else {
      setViewData((prev) => [...prev, ...sliceData[pageNumber]]);
      setIsActive(true);
    }
  }, [pageNumber]);

  useEffect(() => {
    if (!map) return;
    const fetch = async () => {
      const res: AddressInfo = (await getAddress(lat, lng)) as AddressInfo;
      setAddress(res.address_name);
    };
    fetch();
  }, [lat, lng, map]);

  const handleMore = () => {
    if (!sliceData) return;
    if (pageNumber < sliceData.length - 1) setPageNumber((prev) => prev + 1);
    else setPageNumber(0);
  };

  if (curTab === "주변") {
    if (areaLoading || isAreaFetching) {
      return <Skeleton className="w-full h-[40px] rounded-sm bg-grey-dark-1" />;
    }

    if (areaError) {
      return (
        <div>
          <button
            className="absolute right-0 top-2 text-sm text-grey-dark"
            onClick={() => areaRankingRefetch()}
          >
            새로고침
          </button>
          <div className="mb-1">
            <EmojiHoverButton
              className="bg-black-light-2 px-3"
              emoji="🔍"
              text="내 위치"
              subText="내 위치를 기반으로 검색"
              onClickFn={() => {
                centerMapOnCurrentPositionAsync(() => areaRankingRefetch());
              }}
            />
          </div>
          <div className="text-sm text-grey-dark text-center mt-3 mb-2">
            {address} 주변 랭킹
          </div>
          <div className="text-sm">잠시 후 다시 시도해 주세요!</div>
        </div>
      );
    }

    if (areaRanking?.length === 0) {
      return (
        <div>
          <button
            className="absolute right-0 top-2 text-sm text-grey-dark"
            onClick={() => areaRankingRefetch()}
          >
            새로고침
          </button>
          <div className="mb-1">
            <EmojiHoverButton
              className="bg-black-light-2 px-3"
              emoji="🔍"
              text="내 위치"
              subText="내 위치를 기반으로 검색"
              onClickFn={() => {
                centerMapOnCurrentPositionAsync(() => areaRankingRefetch());
              }}
            />
          </div>
          <div className="text-sm text-grey-dark text-center mt-3 mb-2">
            {address} 주변 랭킹
          </div>
          <div className="text-sm">순위에 등록된 위치가 없습니다.</div>
        </div>
      );
    }

    return (
      <div>
        <button
          className="absolute right-0 top-2 text-sm text-grey-dark"
          onClick={() => areaRankingRefetch()}
        >
          새로고침
        </button>
        <div className="mb-1">
          <EmojiHoverButton
            className="bg-black-light-2 px-3"
            emoji="🔍"
            text="내 위치"
            subText="내 위치를 기반으로 검색"
            onClickFn={() => {
              centerMapOnCurrentPositionAsync(() => areaRankingRefetch());
            }}
          />
        </div>
        <div className="text-sm text-grey-dark text-center mt-3 mb-1">
          {address} 주변 랭킹
        </div>

        {areaRanking?.map((marker, index) => {
          return (
            <MarkerListItem
              key={index}
              title={marker.address}
              styleType={"ranking"}
              ranking={index + 1}
              lng={marker.longitude}
              lat={marker.latitude}
              markerId={marker.markerId}
            />
          );
        })}
      </div>
    );
  }

  if (top10?.length === 0) {
    return <div className="text-sm">순위에 등록된 위치가 없습니다.</div>;
  }

  return (
    <div className="">
      {top10?.length === 0 && (
        <div className="text-sm">순위에 등록된 위치가 없습니다.</div>
      )}
      {viewData?.map((marker, index) => {
        return (
          <MarkerListItem
            key={index}
            title={marker.address}
            styleType={"ranking"}
            ranking={index + 1}
            lng={marker.longitude}
            lat={marker.latitude}
            markerId={marker.markerId}
          />
        );
      })}
      {sliceData && viewData.length >= 10 && (
        <div className="relative w-full h-[.4px] bg-grey-dark-1 mt-8 mb-10">
          <button
            className="absolute -top-8 left-1/2 -translate-x-1/2 translate-y-1/2 w-16 h-8 text-sm
        rounded-2xl bg-black border-grey-dark-1 border border-solid"
            onClick={handleMore}
          >
            {pageNumber === sliceData.length - 1 ? "접기" : "더보기"}
          </button>
        </div>
      )}
    </div>
  );
};

export default RankingResult;
