"use client";

import SearchIcon from "@/components/icons/SearchIcon";
import { Skeleton } from "@/components/ui/skeleton";
import useCloseMarkerData from "@/hooks/query/marker/useCloseMarkerData";
import useMapStatusStore from "@/store/useMapStatusStore";
import { useEffect, useRef, useState } from "react";
import MylocateList from "../../mypage/_component/MylocateList";
import ErrorMessage from "@/components/atom/ErrorMessage";

const SearchRangebar = () => {
  const [distance, setDistance] = useState(100);
  const { lat, lng } = useMapStatusStore();

  const {
    data: closeMarker,
    fetchNextPage,
    hasNextPage,
    isLoading,
    isError,
    isFetching,
    refetch,
  } = useCloseMarkerData({ lat, lng, distance });

  const boxRef = useRef(null);

  useEffect(() => {
    const currentRef = boxRef.current;

    const observer = new IntersectionObserver(
      (entries) => {
        const entry = entries[0];
        if (entry.isIntersecting) {
          if (!isFetching && hasNextPage) {
            fetchNextPage();
          }
        }
      },
      { threshold: 0.8 }
    );

    if (currentRef) {
      observer.observe(currentRef);
    }

    return () => {
      if (currentRef) {
        observer.unobserve(currentRef);
      }
    };
  }, [isFetching, hasNextPage, fetchNextPage]);

  return (
    <div className="py-4">
      <div className="flex justify-center items-center gap-2 mb-4">
        <label htmlFor="opacityRange" className="text-white text-sm">
          주변 {distance}m
        </label>
        <input
          type="range"
          id="opacityRange"
          className="range-slider w-[60%] bg-grey-dark-1 rounded-lg appearance-none h-2 cursor-pointer"
          min="100"
          max="5000"
          step="100"
          value={distance}
          onChange={(e) => setDistance(parseInt(e.target.value))}
        />
        <button onClick={() => refetch()}>
          <SearchIcon size={17} />
        </button>
      </div>

      {isLoading && (
        <Skeleton className="w-[90%] mx-auto h-12 rounded-sm bg-black-light-2" />
      )}
      {isError && <ErrorMessage text="잠시 후 다시 시도해 주세요." />}

      <div>
        {closeMarker?.pages[0].markers.length === 0 && (
          <div className="text-center">등록된 장소가 없습니다.</div>
        )}

        <ul>
          {closeMarker?.pages.map((page, index) => {
            return (
              <div key={index}>
                {page.markers.map((marker) => {
                  return (
                    <MylocateList
                      key={marker.markerId}
                      title={marker.address || "지원되지 않는 주소입니다."}
                      subTitle={marker.description || "작성된 설명이 없습니다."}
                      lng={marker.longitude}
                      lat={marker.latitude}
                      markerId={marker.markerId}
                      deleteOption={false}
                    />
                  );
                })}
              </div>
            );
          })}
        </ul>

        {hasNextPage && (
          <div ref={boxRef} className="w-full h-12 px-4">
            <Skeleton className="w-full h-full rounded-sm bg-black-light-2" />
          </div>
        )}
      </div>
    </div>
  );
};

export default SearchRangebar;
