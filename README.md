# chulbong- :kr:
<p align="center">
  <img width="100" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/6236863a-11e1-45d5-b2e2-9bcf40363e1d" alt="k-pullup logo"/></br>
  <img width="1024" alt="2024-09-01 pullupbars" src="https://github.com/user-attachments/assets/eceb488a-8db1-46bb-96c7-1db6bfbcf7a4" />
  <img alt="GitHub commit activity (branch)" src="https://img.shields.io/github/commit-activity/w/Alfex4936/chulbong-kr/main">
</p>

> [!NOTE]  
> The service is running with [Go as backend](https://github.com/Alfex4936/chulbong-kr/tree/main/backend) but [Spring boot 3](https://github.com/Alfex4936/chulbong-kr/tree/main/backend-spring) is also available. (WIP)

<p align="center">
  <img width="356" src="https://github.com/user-attachments/assets/78072f29-59f5-4303-b241-533b9198b18d" alt="gochulbong"/></br>
</p>

### 프로젝트 소개 :world_map:

**chulbong-kr**은 공공장소에 있는 턱걸이 바를 찾고 공유하기 위한 커뮤니티 플랫폼입니다.

지도 API를 활용하여 사용자는 가입 및 로그인 후 턱걸이 바의 위치를 마커로 추가할 수 있으며,

사진 한 장과 간단한 설명을 함께 업로드할 수 있습니다.

로그인한 다른 사용자는 해당 마커에 댓글을 남길 수 있어, 정보 공유 및 소통이 용이합니다.

![0](https://github.com/Alfex4936/chulbong-kr/assets/2356749/c0f58f73-d568-4ef7-8fa1-1f20820b8fff)

|                                                                                                                                                                 |                                                                                                                                                                  |
| :-------------------------------------------------------------------------------------------------------------------------------------------------------------: | :--------------------------------------------------------------------------------------------------------------------------------------------------------------: |
|           <img width="1604" alt="main" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/2ac3ffb2-f22e-4476-bf4b-ebc1bc58b1a4"> 메인 화면            |             <img width="1604" alt="2" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/30967ec3-9921-4910-9c1a-293d950d50ef"> 마커 정보              |
| <img width="1604" alt="road view" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/0a26c86a-6da9-42c8-803a-57d5904dea29"> 거리뷰 (가장 가까운 위치) |             <img width="1604" alt="comment" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/a8d69589-4b7f-435b-8084-ef87419eed09"> 댓글             |
|          <img width="1604" alt="nearby" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/a878595b-d613-4e22-aa95-2d8fee65b578"> 주변 철봉           | <img width="1604" alt="offline-pdf" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/a7eb993e-f847-40c8-9b93-d367f4c6a3f8"> 오프라인 저장 (카카오맵) |
|          <img width="500" alt="consonant" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/d240d3d2-d42c-4136-b483-2805a03231a1"> 철봉 주소 검색 (초성 지원)           | <img width="500" alt="offline-pdf" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/566ac319-4bd9-4acf-8bf4-eb5a05c21b0f"> 정보 수정 제안 |

![slack](https://github.com/Alfex4936/chulbong-kr/assets/2356749/5ec03f6a-871f-4556-90c3-13bb44769f13)

<img width="500" alt="chatting" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/53e4f587-e155-49c7-b28b-d56e150f1fe2">

### 기능

```mermaid
sequenceDiagram
    participant User
    participant AuthService as Auth Service
    participant MarkerService as Marker Service
    participant CommentService as Comment Service
    participant ShareService as Share Service
    participant SearchService as Search Service
    participant AdminService as Admin Service
    participant ChatService as Chat Service
    participant PopularService as Popular Service
    participant OfflineService as Offline Service
    participant AddressService as Address Service

    Note over User, AuthService: User Registration and Authentication
    User->>AuthService: Sign Up with Email
    AuthService->>User: Send Email Verification
    User->>AuthService: Verify Email
    User->>AuthService: Login
    AuthService->>User: Authentication Token

    Note over User, MarkerService: Marker Management
    User->>MarkerService: Add Marker with Location, Photo, Description
    MarkerService->>AddressService: Validate Address
    AddressService-->>MarkerService: Address Validated/Failed
    MarkerService-->>User: Marker Added
    User->>MarkerService: View Marker
    MarkerService-->>User: Display Marker Details

    Note over User, CommentService: Commenting on Markers
    User->>CommentService: Add Comment to Marker
    CommentService-->>User: Comment Added
    User->>CommentService: View Comments
    CommentService-->>User: Display Comments

    Note over User, ShareService: Sharing Markers
    User->>ShareService: Share Marker Link
    ShareService-->>User: Marker Link

    Note over User, SearchService: Searching for Nearby Markers
    User->>SearchService: Search Markers Near Current Location
    SearchService-->>User: Display Nearby Markers

    Note over Admin, AdminService: Admin Tasks
    Admin->>AdminService: Review Markers with Dislikes
    AdminService-->>Admin: Display Markers for Review
    Admin->>AdminService: Update Marker Status
    AdminService-->>Admin: Marker Status Updated

    Note over User, ChatService: Chat Functionality
    User->>ChatService: Join Marker Chat Room
    ChatService-->>User: Joined Chat Room
    User->>ChatService: Send Message in Chat Room
    ChatService-->>User: Message Sent

    Note over User, PopularService: Viewing Popular Markers
    User->>PopularService: View Popular Markers
    PopularService-->>User: Display Popular Markers

    Note over User, OfflineService: Offline Marker Access
    User->>OfflineService: Download Static Map Image
    OfflineService-->>User: Static Map Image with Markers

    Note over User, MarkerService: Suggesting Marker Edits
    User->>MarkerService: Suggest Marker Edit with Photo
    MarkerService-->>User: Edit Suggestion Submitted
```

- **회원가입 및 로그인**: 사용자 인증을 위한 기본적인 회원가입 및 로그인 기능. (이메일 인증 or 소셜 로그인)
- **마커 추가**: 턱걸이 바의 위치를 지도에 마커로 표시. 사진과 간단한 설명을 포함할 수 있음.
- **댓글 기능**: 로그인한 사용자는 각 마커에 댓글을 남길 수 있어 정보 공유가 가능.
- **마커 공유**: 특정 마커 공유 버튼을 눌러서 링크를 공유가 가능.
- **근처 턱걸이 바 검색**: 현재 화면 중앙 위치에서 가까운 턱걸이 바를 찾을 수 있는 기능.
- **관리자**: 자동 1차 필터링 (주소가 없는 경우 db에 기록) + 싫어요 n개 이상 마커들 확인 + 사진 AI 모델 + 추가 제한 구역
- **채팅**: 각 마커마다 채팅 방 + 지역별 채팅 방 (익명)
- **인기 장소 확인**: 사용자들이 실시간 자주 방문하는 인기 턱걸이 바 위치 확인 기능. (현재 위치 기준 + 전국)
- **정적 이미지 오프라인**: 오프라인 용도로 철봉 위치들을 저장할 수 있는 기능. (카카오맵 정적 이미지 보완)
- **마커 장소 검색**: 등록된 마커들의 주소 검색 기능 (초성 지원)
- **정보 수정 제안**: 등록된 마커들의 정보 수정 제안 기능 (사진 1장 필수)

### TODO 아이디어

- **커뮤니티 포럼**: 사용자들이 운동 팁, 턱걸이 바 추천 등을 공유할 수 있는 커뮤니티 공간.
- **이벤트 및 챌린지**: 사용자들이 참여할 수 있는 운동 관련 이벤트 및 챌린지 개최.

### pullup/dips bar detection

"A nimble AI model, trained on 300 images in August 2024, striving to expertly detect bars with precision."

![image](https://github.com/user-attachments/assets/d822d93a-9985-480f-acfc-ba44eb4e96dc)

### 기술 스택

![image](https://github.com/Alfex4936/chulbong-kr/assets/2356749/f82e2295-ce31-4b48-af92-20a8471b7155)

https://github.com/Alfex4936/chulbong-kr/assets/2356749/913b113c-4a8d-4df1-bb5a-83f6babf7475

- **백엔드**: Go언어 Fiber v2, MySQL, AWS S3, LavinMQ (RabbitMQ), Redis, Bleve (Apache Lucene-like, ZincSearch에서 직접 검색 인덱싱으로 변경)
  - 메인: Go, 서브: Java (전체 프로젝트 자바로도 작성 중)
- **프론트엔드**: React -> NextJS (TypeScript)
- **개발 & 운영 효율성**: pprof, flamegraph, Uber's zap logger, Swagger OpenAPI, Prometheus+Grafana
- **AI**: YOLO v8, gpt-4o mini
- **협업**: Slack (+ Slack API)

### 🚀 Project Roles 🚀

- **Backend Development**  
  👨‍💻 [@Alfex4936](https://github.com/Alfex4936)

- **Frontend Development**  
  🎨 [@2YH02](https://github.com/2YH02)

> [!NOTE]
> 대부분 철봉 위치 데이터는 [chulbong.kr](https://chulbong.kr/) 에서 가져왔음을 알립니다. (2차 필터링)
