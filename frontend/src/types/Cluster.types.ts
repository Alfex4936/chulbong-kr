import { type MarkerRes } from "@/api/markers/getAllMarker";
import { type CustomOverlay } from "./CustomOverlay.types";
import { type KakaoMarker } from "./KakaoMap.types";

export interface MarkerClusterer {
  /**
   * 클러스터에 마커 하나를 추가한다.
   *
   * @param marker 추가할 마커
   * @param nodraw 클러스터 redraw 여부. true인 경우 클러스터를 다시 그리지 않는다.
   */
  addMarker(marker: KakaoMarker | CustomOverlay, nodraw?: boolean): void;

  /**
   * 클러스터에 추가된 마커 중 하나를 삭제한다.
   *
   * @param marker 삭제할 마커
   * @param nodraw 클러스터 redraw 여부. true인 경우 클러스터를 다시 그리지 않는다.
   */
  removeMarker(marker: KakaoMarker | CustomOverlay, nodraw?: boolean): void;
  removeMarkers(
    marker: KakaoMarker[] | CustomOverlay | MarkerRes[],
    nodraw?: boolean
  ): void;

  /**
   * 여러개의 마커를 추가한다.
   *
   * @param markers 추가할 마커 객체 배열
   * @param nodraw 클러스터 redraw 여부. true인 경우 클러스터를 다시 그리지 않는다.
   */
  addMarkers(markers: (KakaoMarker | CustomOverlay)[], nodraw?: boolean): void;

  /**
   * 추가된 모든 마커를 삭제한다.
   */
  clear(): void;

  /**
   * 클러스터를 다시 그린다. 주로 옵션을 변경한 이후 클러스터를 다시 그릴 때 사용한다.
   */
  redraw(): void;
}
