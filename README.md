# GuideNow — Realtime Guide-Matching Platform

An enterprise real-time tourism platform that connects travelers with nearby, verified tour guides in Sri Lanka (Colombo, Galle Fort, Kandy, Sigiriya, Weligama, Ella). Built with **Go Microservices**, **Redis Geospatial**, **PostgreSQL**, **WebSocket**, **AI Intent Parsing**, **React/Next.js 15**, **Flutter**, and **AWS**.

---

## 🏛 Platform Architecture

```
                      ┌──────────────────────────────────────┐
                      │    Next.js 15+ Web Explorer HUD     │
                      │   (Interactive Map & AI Prompt Bar)  │
                      └──────────────────┬───────────────────┘
                                         │
                    REST / WebSocket     │   REST / WebSocket (Live GPS)
                                         ▼
                      ┌──────────────────────────────────────┐
                      │         Flutter Mobile App           │
                      │     (Tourist / Guide Dual Mode)      │
                      └──────────────────┬───────────────────┘
                                         │
                                         ▼
                      ┌──────────────────────────────────────┐
                      │     Go API Gateway Service (:8080)   │
                      │  (Reverse Proxy, Rate-Limit, Auth)   │
                      └──────────────────┬───────────────────┘
                                         │
        ┌───────────────┬────────────────┼───────────────┬────────────────┐
        ▼               ▼                ▼               ▼                ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│Guide Service │ │Matching Serv.│ │ Chat Service │ │  AI Service  │ │Booking Serv. │
│   (:8081)    │ │   (:8082)    │ │   (:8083)    │ │   (:8084)    │ │   (:8085)    │
└───────┬──────┘ └───────┬──────┘ └───────┬──────┘ └───────┬──────┘ └───────┬──────┘
        │                │                │                │                │
        ▼                ▼                ▼                ▼                ▼
┌──────────────┐ ┌───────────────────────────────┐ ┌──────────────┐ ┌──────────────┐
│  PostgreSQL  │ │  Redis 7 (GEOSEARCH, Pub/Sub, │ │  Gemini API  │ │  PostgreSQL  │
│(Profiles, DB)│ │     Presence, Rate Limits)    │ │(NLP Fallback)│ │ (Bookings DB)│
└──────────────┘ └───────────────────────────────┘ └──────────────┘ └──────────────┘
```

---

## 🚀 Key Features

### 1. High-Performance Go Microservices
- **API Gateway** (`:8080`): Reverse-proxies all client traffic to backend services with unified CORS, health probes, and route multiplexing.
- **Guide Service** (`:8081`): Manages guide profiles, official SLTDA national licensing credentials, ratings, and live availability toggling.
- **Matching Service** (`:8082`): Sub-10ms nearest-guide geospatial search powered by Redis `GEOSEARCH` and haversine calculations.
- **Chat Service** (`:8083`): Full-duplex WebSocket hub with Redis Pub/Sub multiplexing for multi-room chat, presence, and message persistence.
- **AI Service** (`:8084`): Hybrid NLP query parser translating natural-language tourist prompts into structured database query filters.
- **Booking Service** (`:8085`): State machine coordinating quick-tour requests, real-time guide acceptance, and hourly rate calculations.

### 2. Redis Geospatial Indexing (`GEOSEARCH`)
- Continuous live guide GPS telemetry synchronization into Redis sorted sets (`guides:locations`).
- Rapid radius discovery querying available guides within customizable radii (5 km, 15 km, 25 km) sorted by exact proximity.

### 3. AI-Powered Natural Language Search
- Translates freeform traveler prompts into structured filters:
  - *Query:* `"Looking for an English or German speaking guide in Galle Fort for history and architecture under $25/hr"`
  - *Extracted Filters:*
    ```json
    {
      "location": "Galle",
      "languages": ["English", "German"],
      "specialties": ["Colonial History"],
      "max_rate": 25.0,
      "immediate_only": false
    }
    ```
- Integrates Google Gemini API with automatic fallback to high-accuracy heuristic NLP for offline resilience.

### 4. Interactive Web Explorer HUD (Next.js 15 + Tailwind CSS)
- **Live Radar Canvas**: Animated radar HUD rendering concentric distance rings, coordinate sweeps, and interactive guide blips.
- **AI Prompt Bar**: Instant query parsing with one-click preset prompts and filter badges.
- **WebSocket Chat Drawer**: Live side-over messaging interface with connection indicators.
- **Guide Dispatcher Portal**: Switch between Tourist view and Guide view to toggle live presence and accept incoming tour requests.

### 5. Cross-Platform Mobile Client (Flutter)
- Multi-role application (`flutter/`) supporting both Android and iOS:
  - **Tourist Mode**: Radar map view, quick filters, guide discovery, instant booking dialog.
  - **Guide Mode**: Online/offline status toggle, telemetry GPS pings, incoming tour notification cards.
  - **Live Chat**: Real-time WebSocket messaging interface.

### 6. Cloud-Native & DevOps Ready
- Complete `docker-compose.yml` orchestrating all 6 Go services, PostgreSQL, Redis, and Next.js.
- Production AWS Infrastructure as Code (`infra/`) with Terraform (ECS Fargate, ALB, ElastiCache Redis, RDS PostgreSQL Multi-AZ).
- Automated CI/CD pipeline (`.github/workflows/ci.yml`) executing Go tests, Next.js build, and Flutter test suites.

---

## ⚡ Quickstart Guide

### Option A: Complete Platform with Docker Compose

Start all microservices, databases, and the web frontend with one command:

```bash
docker compose up --build
```

- **Web Explorer HUD**: [http://localhost:3000](http://localhost:3000)
- **API Gateway**: [http://localhost:8080](http://localhost:8080)
- **Guide Service**: [http://localhost:8081](http://localhost:8081)
- **Matching Service**: [http://localhost:8082](http://localhost:8082)
- **Chat Service & WebSocket**: [http://localhost:8083](http://localhost:8083) & `ws://localhost:8080/ws/chat`
- **AI Intent Service**: [http://localhost:8084](http://localhost:8084)
- **Booking Service**: [http://localhost:8085](http://localhost:8085)

---

### Option B: Running Locally

#### 1. Run Go Microservices & Automated Tests

```bash
cd go

# Run all unit and integration tests
go test -v ./...

# Start API Gateway
go run ./cmd/gateway/main.go
```

#### 2. Run Next.js Web Explorer

```bash
cd web
npm install
npm run dev
```
Open [http://localhost:3000](http://localhost:3000) in your browser.

#### 3. Run Flutter Mobile App

```bash
cd flutter
flutter pub get
flutter test
flutter run
```

---

## 📡 API Documentation & Examples

### 1. AI Natural Language Query Parser

```bash
curl -X POST http://localhost:8080/api/v1/ai/parse-query \
  -H "Content-Type: application/json" \
  -d '{"query": "Need a surfing instructor in Weligama under $35 available right now"}'
```

### 2. Geospatial Nearest-Guide Search

```bash
# Query guides within 15km of Galle Fort (6.0329, 80.2168)
curl "http://localhost:8080/api/v1/matching/nearby?lat=6.0329&lng=80.2168&radius=15"
```

### 3. Guide Availability Toggle

```bash
curl -X PATCH http://localhost:8080/api/v1/guides/guide-galle-01/availability \
  -H "Content-Type: application/json" \
  -d '{"available": true}'
```

### 4. Instant Tour Request

```bash
curl -X POST http://localhost:8080/api/v1/bookings/request \
  -H "Content-Type: application/json" \
  -d '{
    "tourist_id": "user-202",
    "tourist_name": "Emma Watson",
    "guide_id": "guide-galle-01",
    "meeting_point": "Galle Fort Lighthouse",
    "destination": "Dutch Heritage Walking Tour",
    "duration_hours": 2.5
  }'
```

### 5. WebSocket Chat Connection

Connect to the WebSocket chat endpoint:
```
ws://localhost:8080/ws/chat?room=room-guide-galle-01&userId=tourist-101&name=Emma&role=tourist
```

Send a message frame:
```json
{"content": "Ayubowan! Are you available for a 2-hour sunset walk?"}
```

---

## ☁️ AWS Infrastructure Deployment (Terraform)

```bash
cd infra
terraform init
terraform plan
terraform apply
```

Provisions:
- Multi-AZ VPC with public and private subnets
- Application Load Balancer with WebSocket upgrade routing (`/ws/*`)
- AWS ECS Fargate cluster for Go microservices
- Amazon RDS PostgreSQL Multi-AZ
- Amazon ElastiCache Redis cluster for geospatial index and Pub/Sub
