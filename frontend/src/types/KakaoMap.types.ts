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
  getCenter: () => LatLngFunctions;
  setLevel: (level: number) => void;
  setCenter: (pos: Pos) => LatLngFunctions;
  panTo: (pos: Pos) => LatLngFunctions;
  getLevel: () => number;
  relayout: VoidFunction;
  addOverlayMapTypeId: (data: any) => void;
  getProjection: () => any;
  setDraggable: (draggable: boolean) => void;
}

export interface KakaoMarker {
  setPosition: (data: Pos & LatLngFunctions) => void;
  getPosition: () => Pos;
  setImage: (img: any) => void;
  setMap: (data: KakaoMap | null | number) => void;
  getTitle: () => string;
  Gb: string;
}

export interface Qa {
  La: number;
  Ma: number;
}

export interface KakaoLatLng {
  center: Qa;
  level: number;
  maxLevel: number;
}
