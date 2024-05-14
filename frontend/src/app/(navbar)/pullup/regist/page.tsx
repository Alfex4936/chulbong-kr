import BlackSideBody from "@/components/atom/BlackSideBody";
import MiniMap from "@/components/map/MiniMap";
import { Separator } from "@/components/ui/separator";
import SearchInput from "../../home/_components/SearchInput";
import Facilities from "./_components/Facilities";
import UploadImage from "./_components/UploadImage";
import MarkerDescription from "./_components/MarkerDescription";
import PrevHeader from "@/components/atom/PrevHeader";

const PullupRegist = () => {
  return (
    <BlackSideBody toggle bodyClass="p-0 mo:px-0 mo:pb-0">
      <PrevHeader url="/home" text="위치 등록" />
      <div className="px-9 mo:px-4">
        <p className="mb-2">🚩 등록 위치를 선택해 주세요</p>
        <SearchInput mini searchToggle />
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
