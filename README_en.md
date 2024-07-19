# chulbong- :kr:
<p align="center">
  <img width="100" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/6236863a-11e1-45d5-b2e2-9bcf40363e1d" alt="k-pullup logo"/></br>
  <img width="1024" alt="2024-07-15 pullupbars" src="https://github.com/user-attachments/assets/01a463c3-6938-4617-9b86-42b0742ba7c3" />
  <img alt="GitHub commit activity (branch)" src="https://img.shields.io/github/commit-activity/w/Alfex4936/chulbong-kr/main">
</p>

> [!NOTE]  
> The service is running with [Go as backend](https://github.com/Alfex4936/chulbong-kr/tree/main/backend) but [Spring boot 3](https://github.com/Alfex4936/chulbong-kr/tree/main/backend-spring) is also available. (WIP)

<p align="center">
  <img width="356" src="https://github.com/user-attachments/assets/78072f29-59f5-4303-b241-533b9198b18d" alt="gochulbong"/></br>
</p>

### Project Introduction :world_map:

**chulbong-kr** is a community platform for finding and sharing pull-up bars in public places.

Using the map API, users can add markers for pull-up bar locations after signing up and logging in, and can upload a photo and a brief description.

Other logged-in users can leave comments on these markers, making it easy to share information and communicate.

![0](https://github.com/Alfex4936/chulbong-kr/assets/2356749/c0f58f73-d568-4ef7-8fa1-1f20820b8fff)

|                                                                                                                                                                 |                                                                                                                                                                  |
| :-------------------------------------------------------------------------------------------------------------------------------------------------------------: | :--------------------------------------------------------------------------------------------------------------------------------------------------------------: |
|           <img width="1604" alt="main" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/2ac3ffb2-f22e-4476-bf4b-ebc1bc58b1a4"> Main Screen           |             <img width="1604" alt="2" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/30967ec3-9921-4910-9c1a-293d950d50ef"> Marker Information     |
| <img width="1604" alt="road view" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/0a26c86a-6da9-42c8-803a-57d5904dea29"> Street View (nearest)      |             <img width="1604" alt="comment" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/a8d69589-4b7f-435b-8084-ef87419eed09"> Comments         |
|          <img width="1604" alt="nearby" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/a878595b-d613-4e22-aa95-2d8fee65b578"> Nearby Pull-up Bars | <img width="1604" alt="offline-pdf" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/a7eb993e-f847-40c8-9b93-d367f4c6a3f8"> Offline Storage (KakaoMap) |
|          <img width="500" alt="consonant" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/d240d3d2-d42c-4136-b483-2805a03231a1"> Pull-up Bar Address Search (Initials Supported)          | <img width="500" alt="offline-pdf" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/566ac319-4bd9-4acf-8bf4-eb5a05c21b0f"> Suggest Information Edit |

![slack](https://github.com/Alfex4936/chulbong-kr/assets/2356749/5ec03f6a-871f-4556-90c3-13bb44769f13)

<img width="500" alt="chatting" src="https://github.com/Alfex4936/chulbong-kr/assets/2356749/53e4f587-e155-49c7-b28b-d56e150f1fe2">

### Features

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

- **Sign Up and Login**: Basic sign-up and login functionality for user authentication. (Email verification required)
- **Add Marker**: Mark the location of pull-up bars on the map. You can include a photo and a brief description.
- **Comment Feature**: Logged-in users can leave comments on each marker to share information.
- **Share Marker**: You can share the link by pressing the share button on a specific marker.
- **Search Nearby Pull-up Bars**: Function to find pull-up bars near the current central position of the screen.
- **Admin**: Automatic first filtering (recorded in the database if no address) + review markers with a certain number of dislikes.
- **Chat**: Chat rooms for each marker + regional chat rooms (anonymous).
- **View Popular Locations**: Function to check the locations of popular pull-up bars that users frequently visit in real-time. (Based on the current location + nationwide)
- **Offline Static Image**: Function to save pull-up bar locations for offline use. (Supplemented with KakaoMap static images)
- **Search Marker Locations**: Function to search for the addresses of registered markers (supports initials).
- **Suggest Information Edit**: Function to suggest edits to the information of registered markers (requires one photo).

### TODO Ideas

- **Community Forum**: A community space where users can share workout tips, recommend pull-up bars, and more.
- **Events and Challenges**: Hosting workout-related events and challenges for users to participate in.

### Technology Stack

![image](https://github.com/Alfex4936/chulbong-kr/assets/2356749/f82e2295-ce31-4b48-af92-20a8471b7155)

https://github.com/Alfex4936/chulbong-kr/assets/2356749/913b113c-4a8d-4df1-bb5a-83f6babf7475

- **Backend**: Go language Fiber v2, MySQL, AWS S3, LavinMQ (RabbitMQ), Redis, Bleve (Apache Lucene-like, replaced with direct search indexing in ZincSearch)
  - Main: Go, Sub: Java (the entire project is also being written in Java)
- **Frontend**: React -> NextJS (TypeScript)
- **Development & Operational Efficiency**: pprof, flamegraph, Uber's zap logger, Swagger OpenAPI, Prometheus+Grafana
- **Collaboration**: Slack (+ Slack API)

### 🚀 Project Roles 🚀

- **Backend Development**  
  👨‍💻 [@Alfex4936](https://github.com/Alfex4936)

- **Frontend Development**  
  🎨 [@2YH02](https://github.com/2YH02)

> [!NOTE]
> Most pull-up bar location data was obtained from [chulbong.kr](https://chulbong.kr/) (second filtering)