"use client";

import MarkerListItem from "@/components/atom/MarkerListItem";
import { Skeleton } from "@/components/ui/skeleton";
import useAreaMarkerRankingData from "@/hooks/query/useAreaMarkerRankingData";
import useMarkerRankingData from "@/hooks/query/useMarkerRankingData";
import useMapStatusStore from "@/store/useMapStatusStore";
import useMapStore from "@/store/useMapStore";
import useTabStore from "@/store/useTabStore";
import getAddress, { type AddressInfo } from "@/utils/getAddress";
import { useEffect, useState } from "react";

const RankingResult = () => {
  const { map } = useMapStore();
  const { lat, lng } = useMapStatusStore();
  const { curTab } = useTabStore();

  const { data: top10 } = useMarkerRankingData();
  const {
    data: areaRanking,
    isError: areaError,
    isLoading: areaLoading,
    refetch: areaRankingRefetch,
    isFetching: isAreaFetching,
  } = useAreaMarkerRankingData(lat, lng, curTab === "주변");

  const [address, setAddress] = useState("");

  useEffect(() => {
    if (!map) return;
    const fetch = async () => {
      const res: AddressInfo = (await getAddress(lat, lng)) as AddressInfo;
      setAddress(res.address_name);
    };
    fetch();
  }, [lat, lng, map]);

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
      {top10?.map((marker, index) => {
        return (
          <MarkerListItem
            key={index}
            title={marker.address}
            styleType={"ranking"}
            ranking={index + 1}
            lng={marker.longitude}
            lat={marker.latitude}
          />
        );
      })}
    </div>
  );
};

export default RankingResult;
