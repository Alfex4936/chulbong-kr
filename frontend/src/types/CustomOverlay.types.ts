import { KakaoMap } from "./KakaoMap.types";
import { Roadview } from "./RoadView.types";

export interface CustomOverlay {
  /**
   * 커스텀 오버레이의 내용을 지정했던 형태로 반환한다.
   */
  getContent(): HTMLElement | string;

  /**
   * 커스텀 오버레이의 표시 여부를 지정한다.
   *
   * @param visible
   */
  setVisible(visible: boolean): void;

  /**
   * 커스텀 오버레이의 표시 여부를 반환한다.
   */
  getVisible(): boolean;

  /**
   * 커스텀 오버레이의 z-index를 설정한다.
   *
   * @param zIndex
   */
  setZIndex(zIndex: number): void;

  /**
   * 커스텀 오버레이의 z-index를 반환한다.
   */
  getZIndex(): number;

  /**
     * 지도 또는 로드뷰에 커스텀 오버레이를 올린다.
     * null 을 지정하면 오버레이를 제거한다.
     *
     * @param map
     */
  setMap(map: KakaoMap | Roadview | null): void;
}