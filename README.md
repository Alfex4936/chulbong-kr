# chulbong-kr

![chul4](https://github.com/Alfex4936/chulbong-kr/assets/2356749/60aba6ca-6339-47a0-b2f2-a9254d722755)
![map](https://github.com/Alfex4936/chulbong-kr/assets/2356749/d1811e5d-9857-4e9f-b997-78b77df343fb)
![slack](https://github.com/Alfex4936/chulbong-kr/assets/2356749/5ec03f6a-871f-4556-90c3-13bb44769f13)

### 프로젝트 
소개
**chulbong-kr**은 공공장소에 있는 턱걸이 바를 찾고 공유하기 위한 커뮤니티 플랫폼입니다.

카카오맵 API를 활용하여 사용자는 가입 및 로그인 후 턱걸이 바의 위치를 마커로 추가할 수 있으며,

사진 한 장과 간단한 설명을 함께 업로드할 수 있습니다.

로그인한 다른 사용자는 해당 마커에 댓글을 남길 수 있어, 정보 공유 및 소통이 용이합니다.

### 기능
- **회원가입 및 로그인**: 사용자 인증을 위한 기본적인 회원가입 및 로그인 기능. (이메일 인증 필요)
- **마커 추가**: 턱걸이 바의 위치를 지도에 마커로 표시. 사진과 간단한 설명을 포함할 수 있음. (Spatial 타입 이용으로 주변 마커 확인)
- **댓글 기능**: 로그인한 사용자는 각 마커에 댓글을 남길 수 있어 정보 공유가 가능.
- **근처 턱걸이 바 검색**: 사용자의 현재 위치에서 가까운 턱걸이 바를 찾을 수 있는 기능.
- **사용자 평가/신고 시스템**: 턱걸이 바 이용 싫어요 기능을 통한 리포트 기능.

### TODO 아이디어
- **커뮤니티 포럼**: 사용자들이 운동 팁, 턱걸이 바 추천 등을 공유할 수 있는 커뮤니티 공간.
- **인기 장소 확인**: 사용자들이 자주 방문하는 인기 턱걸이 바 위치 확인 기능.
- **이벤트 및 챌린지**: 사용자들이 참여할 수 있는 운동 관련 이벤트 및 챌린지 개최.

### 기술 스택
- **백엔드**: Go언어 Fiber v2, AWS RDS MySQL (+ spatial type), AWS S3, AWS EC2 (fly.io), Redis
- **프론트엔드**: React (TypeScript)
- **협업**: Slack (+ Slack API)

### 성능 테스트

환경: AWS RDS (Free tier, MySQL (SRID)) + 24GB RAM PC

100,000개 마커 정보 (20MB) 로딩 약 2초.

![image](https://github.com/Alfex4936/chulbong-kr/assets/2356749/44956afa-8c6c-414f-a6ff-1f11d348c3f5)
