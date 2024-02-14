export type Pos = {
  La: number;
  Ma: number;
};

export type LatLngFunctions = {
  getLat: VoidFunction;
  getLng: VoidFunction;
};

export interface KaKaoMapMouseEvent {
  latLng: Pos & LatLngFunctions;
  point: { x: number; y: number };
}

export interface KakaoMap {
  getCenter: VoidFunction;
  setLevel: (level: number) => void;
  setCenter: (pos: Pos) => void;
}