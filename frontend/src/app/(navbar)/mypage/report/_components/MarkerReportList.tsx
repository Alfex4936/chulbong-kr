import BlackLightBox from "@/components/atom/BlackLightBox";
import GrowBox from "@/components/atom/GrowBox";
import AlertButton from "@/components/common/AlertButton";
import DeleteIcon from "@/components/icons/DeleteIcon";
import { Separator } from "@/components/ui/separator";
import useApproveReport from "@/hooks/mutation/report/useApproveReport";
import useDeleteReport from "@/hooks/mutation/report/useDeleteReport";
import useDenyReport from "@/hooks/mutation/report/useDenyReport";
import useMarkerData from "@/hooks/query/marker/useMarkerData";
import useMapStore from "@/store/useMapStore";
import getAddress, { type AddressInfo } from "@/utils/getAddress";
import Image from "next/image";
import { useEffect, useState } from "react";
import ChangePassword from "../../user/ChangePassword";
import StatusBadge from "./StatusBadge";
import { Skeleton } from "@/components/ui/skeleton";

interface Props {
  markerId: number;
  lat: number;
  lng: number;
  desc: string;
  img: string[];
  status: string;
  userId: number;
  myId?: number;
  reportId: number;
}

interface InfoListProps {
  text: string;
  subText: string;
  buttonText?: string;
  isTruncate?: boolean;
}

const InfoList = ({
  text,
  subText,
  buttonText,
  isTruncate = false,
}: InfoListProps) => {
  return (
    <div className="flex text-[13px] py-1">
      <span className={`w-1/5 text-wrap break-words`}>{text}</span>
      <span
        className={`w-4/5 ${isTruncate ? "truncate" : "text-wrap break-words"}`}
      >
        {subText}
      </span>
      <GrowBox />
      {buttonText && <ChangePassword />}
    </div>
  );
};

const MarkerReportList = ({
  markerId,
  lat,
  lng,
  desc,
  img,
  status,
  userId,
  myId,
  reportId,
}: Props) => {
  const { data: marker, isLoading: markerLoading } = useMarkerData(markerId);
  const { mutate: deleteReport } = useDeleteReport(markerId, reportId);
  const { mutate: approveReport } = useApproveReport(markerId);
  const { mutate: denyReport } = useDenyReport(markerId);
  const { map } = useMapStore();

  const [addr, setAddr] = useState("");
  const [dropdown, setDropdown] = useState(false);

  useEffect(() => {
    if (!map) return;

    const fetchAddr = async () => {
      const data = (await getAddress(lat, lng)) as AddressInfo;
      setAddr(data.address_name);
    };

    fetchAddr();
  }, [map]);

  if (markerLoading)
    return <Skeleton className="w-[90%] p-4 rounded-md h-60 mx-auto" />;

  return (
    <BlackLightBox className="relative">
      <div className="flex items-center mb-2">
        {myId && userId === myId && (
          <AlertButton
            ButtonText={<DeleteIcon size={20} />}
            title="정말 삭제하시겠습니까?"
            clickFn={deleteReport}
          />
        )}
        <GrowBox />
        <div className="relative">
          <button
            onClick={() => {
              if (marker?.isChulbong) setDropdown((prev) => !prev);
            }}
          >
            <StatusBadge status={status} />
          </button>
          {dropdown && status !== "APPROVED" && status !== "DENY" && (
            <div className="absolute top-8 left-0">
              {marker?.isChulbong && (
                <div>
                  <button
                    className="mb-1"
                    onClick={() => approveReport(reportId)}
                  >
                    <StatusBadge status={"APPROVED"} />
                  </button>
                  <button onClick={() => denyReport(reportId)}>
                    <StatusBadge status={"DENY"} />
                  </button>
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      {status !== "APPROVED" && (
        <>
          <div>기존</div>
          <InfoList
            text="주소"
            subText={marker?.address || "제공되는 주소가 없음"}
          />
          <InfoList
            text="설명"
            subText={marker?.description || "작성된 설명 없음"}
            isTruncate
          />
          <Separator className="mx-1 my-3 bg-grey-dark-1" />
        </>
      )}

      <div>수정</div>
      <InfoList text="주소" subText={addr || "제공되는 주소가 없음"} />
      <InfoList text="설명" subText={desc || "작성된 설명 없음"} isTruncate />
      {img && (
        <div>
          <Separator className="mx-1 my-3 bg-grey-dark-1" />
          <div>추가된 이미지</div>
          <div className="flex">
            {img?.map((img) => {
              return (
                <Image
                  src={img as string}
                  width={30}
                  height={30}
                  alt="마커 수정"
                  className="w-10 h-10 object-contain ml-2"
                  key={img}
                />
              );
            })}
          </div>
        </div>
      )}
    </BlackLightBox>
  );
};

export default MarkerReportList;
