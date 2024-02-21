export type Pos = {
  La: number;
  Ma: number;
};

export type LatLngFunctions = {
  getLat: () => number;
  getLng: () => number;
};

export interface KaKaoMapMouseEvent {
  latLng: Pos & LatLngFunctions;
  point: { x: number; y: number };
}

export interface KakaoMap {
  getCenter: VoidFunction;
  setLevel: (level: number) => void;
  setCenter: (pos: Pos) => void;
  getLevel: () => number;
}

export interface KakaoMarker {
  setPosition: (data: Pos & LatLngFunctions) => void;
  setMap: (data: KakaoMap | null | number) => void;
}
