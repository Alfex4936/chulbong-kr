import { type KakaoMap } from "@/types/KakaoMap.types";

class MapWalker {
  walker: any;
  content: HTMLElement;
  startX: any;
  startY: any;
  startOverlayPoint: any;
  map: KakaoMap;
  mapWrapper: HTMLDivElement;
  roadviewContainer: HTMLDivElement;
  roadview: any;
  roadviewClient: any;
  newPos: any;
  mapContainer: any;
  prevPos: any;

  constructor(
    position: any,
    map: KakaoMap,
    mapWrapper: HTMLDivElement,
    roadviewContainer: HTMLDivElement,
    roadview: any,
    roadviewClient: any,
    mapContainer: any
  ) {
    const content: HTMLElement = document.createElement("div");
    const figure: HTMLElement = document.createElement("div");
    const angleBack: HTMLElement = document.createElement("div");

    content.className = "MapWalker";
    figure.className = "figure";
    angleBack.className = "angleBack";

    content.appendChild(angleBack);
    content.appendChild(figure);

    const walker = new window.kakao.maps.CustomOverlay({
      position,
      content,
      yAnchor: 1,
    });

    this.walker = walker;
    this.content = content;
    this.map = map;

    this.mapWrapper = mapWrapper;
    this.mapContainer = mapContainer;
    this.roadviewContainer = roadviewContainer;
    this.roadviewClient = roadviewClient;
    this.roadview = roadview;

    this.newPos;
    this.prevPos;

    this.onMouseMove = this.onMouseMove.bind(this);
    this.onMouseUp = this.onMouseUp.bind(this);
    this.onMouseDown = this.onMouseDown.bind(this);
    this.toggleRoadview = this.toggleRoadview.bind(this);
  }

  setAngle(angle: any) {
    const threshold = 22.5;
    for (let i = 0; i < 16; i++) {
      if (angle > threshold * i && angle < threshold * (i + 1)) {
        const className = "m" + i;
        this.content.className = this.content.className.split(" ")[0];
        this.content.className += " " + className;
        break;
      }
    }
  }

  setPosition(position: any) {
    this.walker.setPosition(position);
    this.newPos = position;
    this.prevPos = position;
  }

  setMap() {
    this.walker.setMap(this.map);
  }

  onMouseDown(e: any) {
    if (e.preventDefault) {
      e.preventDefault();
    } else {
      e.returnValue = false;
    }

    const proj = this.map.getProjection();
    const overlayPos = this.walker.getPosition();

    this.prevPos = overlayPos;

    window.kakao.maps.event.preventMap();

    this.startX = e.clientX;
    this.startY = e.clientY;

    this.startOverlayPoint = proj.containerPointFromCoords(overlayPos);

    this.addEventHandle(this.mapContainer, "mousemove", this.onMouseMove);
  }

  onMouseMove(e: any) {
    if (e.preventDefault) {
      e.preventDefault();
    } else {
      e.returnValue = false;
    }

    const proj = this.map.getProjection();
    const deltaX = this.startX - e.clientX;
    const deltaY = this.startY - e.clientY;

    const newPoint = new window.kakao.maps.Point(
      this.startOverlayPoint.x - deltaX,
      this.startOverlayPoint.y - deltaY
    );
    this.newPos = proj.coordsFromContainerPoint(newPoint);

    this.walker.setPosition(this.newPos);
  }

  onMouseUp() {
    this.removeEventHandle(this.mapContainer, "mousemove", this.onMouseMove);
    this.toggleRoadview();
  }

  addEventHandle(target: any, type: any, callback: any) {
    if (target.addEventListener) {
      target.addEventListener(type, callback);
    } else {
      target.attachEvent("on" + type, callback);
    }
  }

  removeEventHandle(target: any, type: any, callback: any) {
    if (target.removeEventListener) {
      target.removeEventListener(type, callback);
    } else {
      target.detachEvent("on" + type, callback);
    }
  }

  toggleRoadview() {
    if (!this.newPos || this.prevPos === this.newPos) return;

    this.roadviewClient.getNearestPanoId(this.newPos, 50, (panoId: number) => {
      if (panoId === null) {
        this.walker.setPosition(this.prevPos);
      } else {
        this.prevPos = this.newPos;
        this.map.relayout();
        this.walker.setPosition(this.newPos);

        this.roadview.setPanoId(panoId, this.newPos);
        this.roadview.relayout();
      }
    });
  }

  init() {
    this.addEventHandle(this.mapContainer, "mouseup", this.onMouseUp);
    this.addEventHandle(this.content, "mousedown", this.onMouseDown);
  }
}

export default MapWalker;
