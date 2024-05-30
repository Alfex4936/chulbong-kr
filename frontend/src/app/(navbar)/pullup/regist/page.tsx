import BlackSideBody from "@/components/atom/BlackSideBody";
import PrevHeader from "@/components/atom/PrevHeader";
import MapSearch from "@/components/map/MapSearch";
import MiniMap from "@/components/map/MiniMap";
import { Separator } from "@/components/ui/separator";
import Facilities from "./_components/Facilities";
import MarkerDescription from "./_components/MarkerDescription";
import UploadImage from "./_components/UploadImage";

const PullupRegist = () => {
  return (
    <BlackSideBody toggle bodyClass="p-0 mo:px-0 mo:pb-0">
      <PrevHeader url="/home" text="위치 등록" />
      <div className="px-4 pt-2 pb-4 mo:pb-20">
        <p className="mb-2">🚩 등록 위치를 선택해 주세요</p>
        <MapSearch
          mini
          className={`relative top-0 w-full mb-4 bg-black z-[90]`}
        />
        <MiniMap />
        <Separator className="my-4 bg-grey-dark-1" />
        <p className="mb-2">🎁 기구 개수를 입력해 주세요</p>
        <Facilities />
        <Separator className="my-4 bg-grey-dark-1" />
        <p className="mb-2">📷 사진을 등록해 주세요</p>
        <UploadImage />
        <Separator className="my-4 bg-grey-dark-1" />
        <p className="mb-2">📒 설명을 입력해 주세요.</p>
        <MarkerDescription />
      </div>
    </BlackSideBody>
  );
};

export default PullupRegist;
