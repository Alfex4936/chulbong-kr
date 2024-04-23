"use client";

import { Skeleton } from "@/components/ui/skeleton";
import useMylocateData from "@/hooks/query/user/useMylocateData";
import { useEffect, useRef } from "react";
import MylocateList from "../_component/MylocateList";

const MylocateClient = () => {
  const {
    data: mylocates,
    fetchNextPage,
    hasNextPage,
    isLoading,
    isError,
    isFetching,
  } = useMylocateData();

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

  if (isError) return <div>잠시 후 다시 시도해 주세요!</div>;

  if (isLoading) {
    return <Skeleton className="w-full h-12 rounded-sm bg-black-light-2" />;
  }

  return (
    <div>
      {mylocates?.pages[0].markers.length === 0 && (
        <div>등록된 리뷰가 없습니다.</div>
      )}

      <ul>
        {mylocates?.pages.map((page, index) => {
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
  );
};

export default MylocateClient;
