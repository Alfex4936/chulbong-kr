import instance from "@/api/instance";
import SearchInput from "@/app/(navbar)/home/_components/SearchInput";
import BlackSideBody from "@/components/atom/BlackSideBody";
import Heading from "@/components/atom/Heading";
import PrevHeader from "@/components/atom/PrevHeader";
import MiniMap from "@/components/map/MiniMap";
import { Separator } from "@/components/ui/separator";
import { type Marker } from "@/types/Marker.types";
import MarkerDescription from "../../regist/_components/MarkerDescription";
import UploadImage from "../../regist/_components/UploadImage";

const getMarker = async (id: number): Promise<Marker> => {
  const res = await instance.get(
    `${process.env.NEXT_PUBLIC_BASE_URL}/markers/${id}/details`
  );

  return res.data;
};

interface Params {
  id: string;
}

interface Props {
  params: Params;
}

const RportMarkerPage = async ({ params }: Props) => {
  try {
    const marker = await getMarker(Number(params.id));
    return (
      <BlackSideBody toggle bodyClass="relative p-0 mo:px-0 mo:pb-0">
        <PrevHeader
          url={`/pullup/${params.id}/reportlist`}
          text="정보 수정 제안"
        />

        <div className="px-9 pb-4 scrollbar-thin mo:px-4 mo:pb-20">
          <p className="mb-2">🚩 수정할 위치를 선택해 주세요</p>
          <SearchInput mini searchToggle />
          <MiniMap
            isMarker
            latitude={marker.latitude}
            longitude={marker.longitude}
          />
          {/* <Separator className="my-4 bg-grey-dark-1" />
              <p className="mb-2">🎁 기구 개수를 입력해 주세요</p> */}
          {/* <Facilities /> */}
          <Separator className="my-4 bg-grey-dark-1" />
          <p className="mb-2">📷 사진을 등록해 주세요</p>
          <UploadImage />
          <Separator className="my-4 bg-grey-dark-1" />
          <p className="mb-2">📒 설명을 입력해 주세요.</p>
          <MarkerDescription
            desc={marker.description}
            markerId={Number(params.id)}
            isReport={true}
          />
        </div>
      </BlackSideBody>
    );
  } catch (error) {
    return (
      <BlackSideBody toggle>
        <Heading title="정보 수정 제안" />
        <p className="text-center text-red">존재하지 않는 위치입니다.</p>
      </BlackSideBody>
    );
  }
};

export default RportMarkerPage;
