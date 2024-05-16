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
  - [ ] Swagger 연동
  - [ ] JUnit 연동
  - [ ] ElasticSearch or Melisearch 연동
  - [x] JaCoCo/CodeCov 연동
- [ ] controller
  - [ ] marker
    - [x] GET /markers
  - [ ] user
  - [ ] report
  - [ ] comment
  - [ ] ranking
- [ ] service
  - [x] Profanity (욕설 필터링)
    - 벤치마크 완료 (`String.contains` vs `Double-Array Ahocorasick`)
- [ ] repository
- [ ] util
  - [x] TimezoneFinder
- [ ] external APIs