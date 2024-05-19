# chulbong-kr

[![codecov](https://codecov.io/gh/Alfex4936/chulbong-kr/graph/badge.svg?token=R4VBHZKJ7F)](https://codecov.io/gh/Alfex4936/chulbong-kr)

> [!WARNING]
> Work in progress...

chulbong-kr 을 Spring Boot 3을 이용해 만드는 과정

# TODO
- [x] 프로젝트 기본 설정
  - [x] Security 6
  - [x] MySQL 연동
  - [x] JPA 연동
  - [x] Lombok 연동
  - [x] Config 연동
  - [x] JUnit 연동
  - [x] JaCoCo/CodeCov 연동
  - [x] openjdk JMH 연동
  - [ ] Swagger 연동
  - [ ] Zincsearch or Meilisearch 연동
- [ ] controller
  - [ ] marker
    - [x] GET /markers
    - [x] GET /markers/close
    - Location
      - [x] GET /markers/save-offline
      - [x] GET /markers/weather
  - [ ] user
  - [ ] report
  - [ ] comment
  - [ ] ranking
- [ ] service
  - [x] Profanity (욕설 필터링)
    - 벤치마크 완료 (`String.contains` vs `Double-Array Ahocorasick`)
  - [x] PDF (오프라인 PDF 생성)
    - [x] File Download (파일 다운로드 from URL)
    - [x] Image 처리
    - [x] Temp folder cleanup scheduler (임시 폴더 청소)
  - [ ] marker
- [ ] repository
  - [x] Marker
- [ ] util
  - [x] TimezoneFinder to check South Korea
  - [x] WCONGNAMUL to WGS84
  - [x] WGS84 to WCONGNAMUL
- [ ] external APIs
  - [x] Kakao Map Address API