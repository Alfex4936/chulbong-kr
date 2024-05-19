# chulbong- :kr:

> [!NOTE]  
> The service is running with [Go as backend](https://github.com/Alfex4936/chulbong-kr/tree/main/backend) but [Spring boot 3](https://github.com/Alfex4936/chulbong-kr/tree/main/backend-spring) is also available. (WIP)

### 프로젝트 소개 :world_map:
**chulbong-kr**은 공공장소에 있는 턱걸이 바를 찾고 공유하기 위한 커뮤니티 플랫폼입니다.

지도 API를 활용하여 사용자는 가입 및 로그인 후 턱걸이 바의 위치를 마커로 추가할 수 있으며,

사진 한 장과 간단한 설명을 함께 업로드할 수 있습니다.

로그인한 다른 사용자는 해당 마커에 댓글을 남길 수 있어, 정보 공유 및 소통이 용이합니다.

![0](https://github.com/Alfex4936/chulbong-kr/assets/2356749/c0f58f73-d568-4ef7-8fa1-1f20820b8fff)

| | |
|:-------------------------:|:-------------------------:|
| <img width="1604" alt="1" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/6e45dac1-6c4a-4b84-bf47-10e8c09f4f2b"> 메인 화면 | <img width="1604" alt="2" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/2bcfe0f8-eb55-46d8-ba17-0274f25bfa38"> 마커 정보 (댓글, 공유, 거리뷰)|
| <img width="1604" alt="3" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/d56a0e34-22fa-4124-bafc-f3b1c0a3f90c"> 거리뷰 (가장 가까운 위치)| <img width="1604" alt="4" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/7c2f5a39-a82f-47b0-8cfa-14bfff0551bd"> 댓글|

![chatting](https://github.com/Alfex4936/chulbong-kr/assets/2356749/53e4f587-e155-49c7-b28b-d56e150f1fe2)
![slack](https://github.com/Alfex4936/chulbong-kr/assets/2356749/5ec03f6a-871f-4556-90c3-13bb44769f13)
![image](https://github.com/Alfex4936/chulbong-kr/assets/2356749/a7eb993e-f847-40c8-9b93-d367f4c6a3f8)


### 기능
- **회원가입 및 로그인**: 사용자 인증을 위한 기본적인 회원가입 및 로그인 기능. (이메일 인증 필요)
- **마커 추가**: 턱걸이 바의 위치를 지도에 마커로 표시. 사진과 간단한 설명을 포함할 수 있음.
- **댓글 기능**: 로그인한 사용자는 각 마커에 댓글을 남길 수 있어 정보 공유가 가능.
- **마커 공유**: 특정 마커 공유 버튼을 눌러서 링크를 공유가 가능.
- **근처 턱걸이 바 검색**: 현재 화면 중앙 위치에서 가까운 턱걸이 바를 찾을 수 있는 기능.
- **관리자**: 자동 1차 필터링 (주소가 없는 경우 db에 기록) + 싫어요 n개 이상 마커들 확인
- **채팅**: 각 마커마다 채팅 방 + 지역별 채팅 방 (익명)
- **인기 장소 확인**: 사용자들이 실시간 자주 방문하는 인기 턱걸이 바 위치 확인 기능. (현재 위치 기준 + 전국)
- **정적 이미지 오프라인**: 오프라인 용도로 철봉 위치들을 저장할 수 있는 기능. (카카오맵 정적 이미지 보완)
- **마커 장소 검색**: 등록된 마커들의 주소 검색 기능.

### TODO 아이디어
- **커뮤니티 포럼**: 사용자들이 운동 팁, 턱걸이 바 추천 등을 공유할 수 있는 커뮤니티 공간.
- **이벤트 및 챌린지**: 사용자들이 참여할 수 있는 운동 관련 이벤트 및 챌린지 개최.


### 기술 스택
![image](https://github.com/Alfex4936/chulbong-kr/assets/2356749/f82e2295-ce31-4b48-af92-20a8471b7155)


- **백엔드**: Go언어 Fiber v2, MySQL, AWS S3, LavinMQ (RabbitMQ), Redis, ZincSearch (ElasticSearch)
  - 메인: Go, 서브: Java (전체 프로젝트 자바로도 작성 중)
- **프론트엔드**: React (TypeScript)
- **개발 & 운영 효율성**: pprof, flamegraph, Uber's zap logger, Swagger OpenAPI, Prometheus+Grafana
- **협업**: Slack (+ Slack API)
