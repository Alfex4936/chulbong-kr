import BookmarkIcon from "@/components/icons/BookmarkIcon";
import IconButton from "./_components/IconButton";
import ShareIcon from "@/components/icons/ShareIcon";
import DislikeIcon from "@/components/icons/DislikeIcon";
import DeleteIcon from "@/components/icons/DeleteIcon";
import { Separator } from "@/components/ui/separator";
import RoadViewIcon from "@/components/icons/RoadViewIcon";
import Link from "next/link";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import ImageList from "./_components/ImageList";

interface Props {
  markerId: number;
}

const PullupClient = ({ markerId }: Props) => {
  return (
    <div>
      {/* 이미지 배경 */}
      <div
        className="relative w-full h-64 bg-cover bg-center"
        style={{
          backgroundImage: "url('/metaimg.webp')",
        }}
      >
        <div className="absolute top-1 left-1 flex  items-center py-1 px-2 rounded-sm z-20 bg-black-light-2">
          <img
            className="mr-2"
            src="https://t1.daumcdn.net/localimg/localimages/07/2018/pc/weather/ico_weather20.png"
            alt=""
          />
          <span className="text-lg font-bold">15.9℃</span>
        </div>
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
              <span>평행봉</span>
              <span>3개</span>
            </div>
            <Separator className="my-2 bg-grey-dark" />
            <div className="flex justify-between">
              <span>평행봉</span>
              <span>3개</span>
            </div>
          </div>
        </div>
        {/* 정보 */}
        <div className="mt-4">
          <div className="truncate flex items-center mb-[2px]">
            <span className="mr-1">
              <h1 className="text-xl">서울 종로구</h1>
            </span>
            <button>
              <RoadViewIcon />
            </button>
          </div>

          <div className="text-xs text-gray-400 mb-5">
            <span>...등록</span>
            <span>(...업데이트)</span>
            <span className="mx-1">|</span>
            <Link href={"/pullup/3000"} className="underline">
              정보 수정 제안
            </Link>
          </div>

          <h2 className="">작성된 설명이 없습니다.</h2>
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
            <ImageList />
          </TabsContent>
          <TabsContent value="review">Change your password here.</TabsContent>
        </Tabs>
      </div>
    </div>
  );
};

export default PullupClient;
