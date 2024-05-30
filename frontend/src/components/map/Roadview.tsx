import useRoadviewStatusStore from "@/store/useRoadviewStatusStore";
import { type KakaoMap } from "@/types/KakaoMap.types";
import { useEffect, useRef, useState } from "react";
import MapWalker from "./MapWalker";
import { useToast } from "../ui/use-toast";

const Roadview = () => {
  const { lat, lng, close } = useRoadviewStatusStore();

  const { toast } = useToast();

  const [mapHover, setMapHover] = useState(false);
  const [mapData, setMapData] = useState<KakaoMap | null>(null);

  const mapContainer = useRef<HTMLDivElement>(null);
  const mapWrapper = useRef<HTMLDivElement>(null);
  const roadviewContainer = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!mapContainer.current || !roadviewContainer.current || mapData) return;

    const mapCenter = new window.kakao.maps.LatLng(lat, lng);
    const mapOption = {
      center: mapCenter,
      level: 3,
    };

    const map = new window.kakao.maps.Map(mapContainer.current, mapOption);
    setMapData(map);

    const imageSize = new window.kakao.maps.Size(39, 39);
    const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

    const activeMarkerImg = new window.kakao.maps.MarkerImage(
      "/activeMarker.svg",
      imageSize,
      imageOption
    );

    const marker = new window.kakao.maps.Marker({
      position: new window.kakao.maps.LatLng(lat, lng),
      image: activeMarkerImg,
    });

    marker.setMap(map);

    const roadview = new window.kakao.maps.Roadview(roadviewContainer.current);
    const roadviewClient = new window.kakao.maps.RoadviewClient();

    roadviewClient.getNearestPanoId(mapCenter, 50, (panoId: number) => {
      if (panoId === null) {
        toast({ description: "로드뷰를 지원하지 않는 주소입니다." });
        close();
      } else {
        roadview.setPanoId(panoId, mapCenter);
      }
    });

    let mapWalker: any = null;

    window.kakao.maps.event.addListener(roadview, "init", () => {
      // 로드뷰에 마커 표시
      const rMarker = new window.kakao.maps.Marker({
        position: mapCenter,
        map: roadview,
      });

      const projection = roadview.getProjection();

      const viewpoint = projection.viewpointFromCoords(
        rMarker.getPosition(),
        rMarker.getAltitude()
      );
      roadview.setViewpoint(viewpoint);

      // 맵 워커 생성, 상태 변경
      mapWalker = new MapWalker(
        mapCenter,
        map,
        mapWrapper.current as HTMLDivElement,
        roadviewContainer.current as HTMLDivElement,
        roadview,
        roadviewClient,
        mapContainer.current as HTMLDivElement
      );
      mapWalker.setMap();
      mapWalker.init();

      window.kakao.maps.event.addListener(roadview, "viewpoint_changed", () => {
        const viewpoint = roadview.getViewpoint();
        mapWalker.setAngle(viewpoint.pan);
      });

      window.kakao.maps.event.addListener(roadview, "position_changed", () => {
        const position = roadview.getPosition();
        mapWalker.setPosition(position);
        map.setCenter(position);
      });
    });
  }, []);

  useEffect(() => {
    if (!mapData) return;

    if (mapHover)
      mapData.addOverlayMapTypeId(window.kakao.maps.MapTypeId.ROADVIEW);
    else mapData.addOverlayMapTypeId(window.kakao.maps.MapTypeId.ROADMAP);
  }, [mapHover]);

  return (
    <div ref={mapWrapper}>
      <div
        ref={roadviewContainer}
        className="absolute top-0 left-0 w-full h-full z-50 mo:h-[60%]"
      />
      <div
        ref={mapContainer}
        className="absolute bottom-5 left-5 z-50 w-80 h-52 rounded-sm mo:h-[calc(40%-56px)] mo:w-full mo:top-auto mo:bottom-14 mo:left-0"
        onMouseEnter={() => setMapHover(true)}
        onMouseLeave={() => setMapHover(false)}
      />
      <button
        className="absolute top-5 right-5 z-50 px-2 py-1 rounded-sm bg-[rgba(0,0,0,0.7)] hover:bg-black-tp-dark"
        onClick={close}
      >
        닫기
      </button>
    </div>
  );
};

export default Roadview;
