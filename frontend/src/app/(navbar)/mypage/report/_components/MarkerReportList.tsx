import BlackLightBox from "@/components/atom/BlackLightBox";
import GrowBox from "@/components/atom/GrowBox";
import { Separator } from "@/components/ui/separator";
import useMarkerData from "@/hooks/query/marker/useMarkerData";
import useMapStore from "@/store/useMapStore";
import getAddress, { type AddressInfo } from "@/utils/getAddress";
import Image from "next/image";
import { useEffect, useState } from "react";
import ChangePassword from "../../user/ChangePassword";

interface Props {
  markerId: number;
  lat: number;
  lng: number;
  desc: string;
  img?: string | string[];
}

interface InfoListProps {
  text: string;
  subText: string;
  buttonText?: string;
}

const InfoList = ({ text, subText, buttonText }: InfoListProps) => {
  return (
    <div className="flex text-[13px] py-1">
      <span className="w-1/5">{text}</span>
      <span className="w-4/5">{subText}</span>
      <GrowBox />
      {buttonText && <ChangePassword />}
    </div>
  );
};

const MarkerReportList = ({ markerId, lat, lng, desc, img }: Props) => {
  const { data: marker } = useMarkerData(markerId);
  const { map } = useMapStore();
  const [addr, setAddr] = useState("");

  useEffect(() => {
    if (!map) return;

    const fetchAddr = async () => {
      const data = (await getAddress(lat, lng)) as AddressInfo;
      setAddr(data.address_name);
    };

    fetchAddr();
  }, [map]);

  return (
    <BlackLightBox>
      <div>기존</div>
      <InfoList
        text="주소"
        subText={marker?.address || "제공되는 주소가 없음"}
      />
      <InfoList
        text="설명"
        subText={marker?.description || "작성된 설명 없음"}
      />
      <Separator className="mx-1 my-3 bg-grey-dark-1" />
      <div>수정</div>
      <InfoList text="주소" subText={addr || "제공되는 주소가 없음"} />
      <InfoList text="설명" subText={desc || "작성된 설명 없음"} />
      {img && (
        <div>
          <Separator className="mx-1 my-3 bg-grey-dark-1" />
          <div>추가된 이미지</div>
          <div className="flex">
            {typeof img === "object" ? (
              img.map((img) => {
                return (
                  <Image
                    src={img as string}
                    width={30}
                    height={30}
                    alt="마커 수정"
                    className="w-10 h-10 object-contain"
                  />
                );
              })
            ) : (
              <Image
                src={img as string}
                width={30}
                height={30}
                alt="마커 수정"
                className="w-10 h-10 object-contain"
              />
            )}
          </div>
        </div>
      )}
    </BlackLightBox>
  );
};

export default MarkerReportList;
