"use client";

import { FacilitiesRes } from "@/api/markers/getFacilities";
import BookmarkIcon from "@/components/icons/BookmarkIcon";
import DeleteIcon from "@/components/icons/DeleteIcon";
import DislikeIcon from "@/components/icons/DislikeIcon";
import RoadViewIcon from "@/components/icons/RoadViewIcon";
import ShareIcon from "@/components/icons/ShareIcon";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import useFacilitiesData from "@/hooks/query/marker/useFacilitiesData";
import useMarkerData from "@/hooks/query/marker/useMarkerData";
import useWeatherData from "@/hooks/query/marker/useWeatherData";
import formatDate from "@/utils/formatDate";
import formatFacilities from "@/utils/formatFacilities";
import Link from "next/link";
import { useMemo } from "react";
import IconButton from "./_components/IconButton";
import ImageList from "./_components/ImageList";
import ReviewList from "./_components/ReviewList";

// https://local.k-pullup.com:5173/pullup/5329

interface Props {
  markerId: number;
}

const PullupClient = ({ markerId }: Props) => {
  const { data: marker, isError } = useMarkerData(markerId);
  const { data: facilities } = useFacilitiesData(markerId);

  const { data: weather, isLoading: weatherLoading } = useWeatherData(
    marker?.latitude as number,
    marker?.longitude as number,
    !!marker
  );

  const facilitiesData = useMemo(() => {
    return formatFacilities(facilities as FacilitiesRes[]);
  }, [facilities]);

  if (isError) return <div>X</div>;
  if (!marker) return;

  return (
    <div>
      {/* 이미지 배경 */}
      <div
        className="relative w-full h-64 bg-cover bg-center"
        style={{
          backgroundImage: marker.photos
            ? `url(${marker.photos[0].photoUrl})`
            : "url('/metaimg.webp')",
        }}
      >
        {weatherLoading ? (
          <Skeleton className="absolute top-1 left-1 w-28 h-12 bg-black-light-2" />
        ) : (
          <div className="absolute top-1 left-1 flex  items-center py-1 px-2 rounded-sm z-20 bg-black-light-2">
            <img
              className="mr-2"
              src={weather?.iconImage}
              alt={weather?.desc}
            />
            <span className="text-lg font-bold">{weather?.temperature}℃</span>
          </div>
        )}

        <IconButton right={10} top={10} icon={<BookmarkIcon />} />
        <IconButton right={10} top={50} icon={<ShareIcon />} />
        <IconButton right={10} top={90} icon={<DislikeIcon />} />
        <IconButton right={10} top={130} icon={<DeleteIcon />} />
        <div className="absolute top-0 left-0 w-full h-full bg-black-tp-dark z-10" />
      </div>
      {/* 기구 숫자 카드 */}
      <div className="relative z-30 px-9 -translate-y-14 mo:px-4">
        <div className="h-28">
          <div
            className="bg-black-light-2 flex flex-col justify-center mx-auto 
                        h-full shadow-md w-2/3 py-4 px-10 rounded-sm mo:text-sm mo:px-5 mo:py-2"
          >
            <div className="flex justify-between">
              <span>철봉</span>
              <span>{facilitiesData.철봉}개</span>
            </div>
            <Separator className="my-2 bg-grey-dark" />
            <div className="flex justify-between">
              <span>평행봉</span>
              <span>{facilitiesData.평행봉}개</span>
            </div>
          </div>
        </div>
        {/* 정보 */}
        <div className="mt-4">
          <div className="flex items-center mb-[2px]">
            <span className="mr-1">
              <h1 className="whitespace-normal overflow-visible break-words w-3/4 truncate text-xl">
                {marker.address || "제공되는 주소가 없습니다."}
              </h1>
            </span>
            <button>
              <RoadViewIcon />
            </button>
          </div>

          <div className="text-xs text-gray-400 mb-5">
            <span>{formatDate(marker.createdAt)}</span>
            <span>({formatDate(marker.updatedAt)}업데이트)</span>
            <span className="mx-1">|</span>
            <Link href={"/pullup/3000"} className="underline">
              정보 수정 제안
            </Link>
          </div>

          <h2>{marker.description || "작성된 설명이 없습니다."}</h2>
        </div>

        <Separator className="my-3 bg-grey-dark" />

        <Tabs defaultValue="photo" className="w-full">
          <TabsList className="w-full">
            <TabsTrigger className="w-1/2" value="photo">
              사진
            </TabsTrigger>
            <TabsTrigger className="w-1/2" value="review">
              리뷰
            </TabsTrigger>
          </TabsList>
          <TabsContent value="photo">
            <ImageList photos={marker.photos} />
          </TabsContent>
          <TabsContent value="review">
            <ReviewList markerId={marker.markerId} />
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
};

export default PullupClient;
