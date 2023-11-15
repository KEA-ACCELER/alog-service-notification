 # 🌼 동시편집 기능을 제공하는 릴리즈 노트 공유 시스템, A-LOG

**개발 기간** 2023.06 ~ 2023.08 <br/>
**사이트 바로가기** https://alog.acceler.kr/ (🔧업데이트 중) <br/>
**Team repo** https://github.com/orgs/KEA-ACCELER/repositories <br/>

# 🐳 Overview Architecture

![image](https://github.com/KEA-ACCELER/alog-service-project/assets/80394866/b9f31a1a-6375-4f6e-af24-02d4b308002a)

Here, the domains(service)'s relationship is shown. <br/>

![image](https://github.com/KEA-ACCELER/alog-service-project/assets/80394866/c639bc22-3a8f-4b6d-ac4f-99f4004a38bb)

# 📚  Implementation

The 'Notification' service receives messages from other services, stores them in the ScyllaDB database, and delivers notifications upon request from the front end, allowing users to receive various events and updates from applications in real time.

## Service Flow and Features

### 1. Receive and store messages

- Messages are sent to the 'Notification' service from other services using OpenFeign.
- The received messages are stored in ScyllaDB to be retrieved and delivered later.

### 2. Deliver notification messages

- If a notification message is required by the front end, a request is sent to the service endpoint.
- The 'Notification' service responds to the request by retrieving the stored messages from ScyllaDB.


## Interface

The 'Notification' service provides the following main endpoints:

- POST **`/api/noti`**: Receives messages from other services and stores them in ScyllaDB.
- GET **`/api/noti`**: If a notification message is requested by the front end, the stored messages are returned.
- PATCH `/api/noti` : If the notification is confirmed by the front end, it is stored as true.

## Dependencies
```
require (
	github.com/gocql/gocql v1.5.2
	github.com/gofiber/fiber/v2 v2.48.0
	github.com/google/uuid v1.3.0
)
```
## ERD
![image](https://github.com/KEA-ACCELER/alog-service-project/assets/80394866/9450963c-6df9-45aa-9bfb-d6bb6e2d9baf)



# ✨ Installation

## Running the user app only 

- use docker-compose.yml
```
docker compose up -d
```

# 📝 Conclusion and Suggestion

## **If try to improve the quality**

- **Code optimization**: Focus on improving response speed by making full use of Fiber's asynchronous processing function.
- **Add unit tests**: Write unit tests for each function to ensure code stability and reliability.

## **If try to improve the performance**

- **Use asynchronous pattern**: Make full use of Fiber's asynchronous function to increase response speed through parallel processing.
- **Database optimization**: Improve database performance through ScyllaDB indexing and query optimization.

## **Conclusion**

The 'Notification' service uses Fiber of Go language to implement message processing and notification message delivery function. Users can receive various notifications in real time, and experience efficient processing through Fiber's lightweight threads and ScyllaDB performance. In the future, we will gather user feedback and continue to develop and maintain the service through updates.

