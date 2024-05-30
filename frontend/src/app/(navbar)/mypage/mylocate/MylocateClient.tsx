"use client";

import { Skeleton } from "@/components/ui/skeleton";
import useMylocateData from "@/hooks/query/user/useMylocateData";
import { useEffect, useRef } from "react";
import LinkEmojiButton from "../../home/_components/LinkEmojiButton";
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
      <LinkEmojiButton
        url="/pullup/register"
        text="위치 등록"
        subText="위치를 등록하고 다른 사람들과 공유하세요!"
        emoji="🚩"
      />
      {mylocates?.pages[0].markers.length === 0 && (
        <div className="text-center">등록된 장소가 없습니다.</div>
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
                    isFetching={isFetching}
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
