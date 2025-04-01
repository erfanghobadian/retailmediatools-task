# Software Engineering Task: Ad Bidding Service


## Implementation Status


## Assumptions

- There is no huge delay between bidding and impression sending otherwise it might get overspent
- it is a hiring task and performance is not at the best level otherwise there are way more better data fetching using redis or ...
- There is no feature for analytics of events 
- Overengineering is prohibited in the hiring task 
- We can have later discussion about future optimisations and issues with my approach



##  Features & Improvements

### 1. Line Item Management
**Completed:**
- API endpoint for creating line items with comprehensive validation
- Added repository style for db to make change easy
- Added validation rules for all inputs
- Unit test of handler logic
- Using postgres with some indexes

**Future Improvements:**
- Add More db indexes
- Introduce admin controls for managing line item lifecycle (pause, archive)

---

### 2. Tracking System
**Completed:**
- Endpoint for recording ad interaction events (impression, click, conversion)
- Event normalization and validation
- Test-friendly mock repository implementation
- Event storage structured for future analytics

**Future Improvements:**
- Store events in ClickHouse for real-time analytical queries
- Use Redis to cache/save aggregated counts for fast scoring
- Add support for batch event ingestion via Kafka
- Add rate-limiting and authentication to prevent abuse

---

### 3. Bidding System
**Completed:**
- Dynamic bid estimation per ad based on real-time conversion data
- Multi-level performance analysis: item, placement, and global scope
- Built-in pacing logic to prevent early budget exhaustion

**Future Improvements:**
- Add overspending prevention (budget reservation)
- Support for predictive bidding based on ML models
- Make pacing configurable per advertiser or campaign

---

### 4. Bid Scoring Strategy
**Completed:**
- Modular scoring system via the Strategy pattern
- Normalized score-to-bid mapping with min/max bounds

**Future Improvements:**
- Add additional strategies (CTR, hybrid, ML-based)
- Integrate A/B testing to compare strategy effectiveness
- Support weighted scoring combining multiple signals
- Add Scoring based on Category and Keyword 

---

### 5. Budget Tracking
**Completed:**
- Tracks daily spending (CPM-based) for each line item
- Deducts budget on impression events
- Daily reset via scheduled job (cron)

**Future Improvements:**
- Adding db transaction for event recording and spending increase
- Audit logs of spending per line item
- Metrics endpoint for budget trends and burn rate

---

### 6. Extensibility & Architecture
**Completed:**
- Clean separation of concerns via services and interfaces
- Strategy-based scoring allows for easy integration of CTR/ML-based models
- Designed for plug-and-play repository implementations (e.g., SQL, NoSQL, in-memory)

**Future Improvements:**
- Modularize components for deployment as microservices

---

### 7. Testing Infrastructure
**Completed:**
- Isolated unit tests for handler
- Reusable test utilities and fixtures
- Mocked repositories for clean service-level testing

**Future Improvements:**
- Add strategy and pacing tests
- Full end-to-end tests for ad selection flow
- Use CI workflows to enforce test coverage thresholds

---

### 8. Monitoring, Observability, and Security
**Future Work (Planned Across Features):**
- Prometheus metrics and Grafana dashboards
- Structured JSON logs with trace IDs
- Alerting on pacing budget overruns
- Authentication and advertiser-specific rate limiting
- Fraud detection using anomaly-based rules



---
## Overview

You are tasked with extending an ad bidding service responsible for managing and serving advertisements.

The service is built using Go with the Fiber framework and provides basic functionality for creating ad line items.

## Task Requirements

Your challenge is to:

1. **Implement ad selection logic** to find winning ads based on various criteria
2. **Implement a tracking endpoint** to record user interactions with ads
3. **Add relevancy scoring** to improve ad matching quality
4. **Implement appropriate validation** for all endpoints
5. **Document your approach** and any assumptions made

## Prerequisites

* [Go](https://golang.org/doc/install) 1.24+
* [Docker](https://docs.docker.com/engine/install/)
* [Compose](https://docs.docker.com/compose/install/)

## Setup & Environment

This repository provides a basic service structure to get started:

```bash
# Build and start the service
docker-compose up -d

# Check service status
curl http://localhost:8080/health

# Test creating a line item
curl -X POST http://localhost:8080/api/v1/lineitems \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Summer Sale Banner",
    "advertiser_id": "adv123",
    "bid": 2.5,
    "budget": 1000.0,
    "placement": "homepage_top",
    "categories": ["electronics", "sale"],
    "keywords": ["summer", "discount"]
  }'

# Get winning ads for a placement (you'll need to implement this)
curl -X GET "http://localhost:8080/api/v1/ads?placement=homepage_top&category=electronics&keyword=discount"
```

## Configuration

The service uses environment variables for configuration, using [Kelsey Hightower's envconfig](https://github.com/kelseyhightower/envconfig) library.

Available environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| APP_NAME | Application name | "Ad Bidding Service" |
| APP_ENVIRONMENT | Running environment | "development" |
| APP_LOG_LEVEL | Log level (debug, info, warn, error) | "info" |
| APP_VERSION | Application version | "1.0.0" |
| SERVER_PORT | HTTP server port | 8080 |
| SERVER_TIMEOUT | Server timeout for requests | "30s" |

## API Structure

The service exposes the following endpoints:

- **POST /api/v1/lineitems**: Create new ad line items with bidding parameters
- **GET /api/v1/ads**: Get winning ads for a specific placement with optional filters (you'll need to implement this)
- **POST /api/v1/tracking**: Record ad interactions (you'll need to implement this)

The complete API specification is available in the OpenAPI document at `api/openapi.yaml`.

## Data Model

The core data model includes:

- **LineItem**: An advertisement with associated bid information
  - `id`: Unique identifier
  - `name`: Display name of the line item
  - `advertiser_id`: ID of the advertiser
  - `bid`: Maximum bid amount (CPM)
  - `budget`: Daily budget for the line item
  - `placement`: Target placement identifier
  - `categories`: List of associated categories
  - `keywords`: List of associated keywords

## Deliverables

Please provide the following:

1. **Ad Selection Logic**: Implement the logic to select winning ads based on placement, categories, and keywords
2. **Tracking Endpoint**: Implement an endpoint to record impressions, clicks, and conversions
3. **Relevancy System**: Develop a scoring mechanism to determine ad relevance
4. **Input Validation**: Add appropriate validation for all API endpoints
5. **Documentation**: Update the README and API docs with your changes

## Evaluation Criteria

Your solution will be evaluated based on:

- **Code quality**: Clean, well-structured, and maintainable code
- **API design**: RESTful design, appropriate error handling, and documentation
- **Implementation quality**: Performance, reliability, and adherence to Go best practices
- **Documentation**: Clear explanation of your approach, design decisions, and trade-offs
- **Testing**: Comprehensive test coverage and consideration of edge cases
- **Innovation**: Creative solutions to the technical challenges presented

## Technical Requirements

- Your solution should be containerized and runnable with docker-compose
- All code should follow Go best practices and conventions
- The API should handle appropriate error cases with meaningful status codes and messages
- Your implementation should consider performance and scaling aspects
- Update the OpenAPI specification to match your implementation

## Storage Solutions

The current implementation uses in-memory storage for simplicity, but this is not suitable for production. You are free to use any storage solution you prefer.
Choose solutions that best fit the requirements and consider factors like scalability, reliability, and performance.

## Scaling Considerations

As part of your solution, please include a section in your documentation addressing the following questions:

1. How would you scale this service to handle millions of ad requests per minute?
2. What bottlenecks do you anticipate and how would you address them?
3. How would you design the system to ensure high availability and fault tolerance?
4. What data storage and access patterns would you recommend for different components (line items, tracking events, etc.)?
5. How would you implement caching to improve performance?

## Getting Started

1. Clone this repository
2. Explore the existing code to understand the current implementation
3. Run the service locally using docker-compose
4. Implement the required features
5. Update tests and documentation
6. Submit your solution

For local development:
- Build and run: `go run ./cmd/server`
- Run tests: `go test ./...`
- Build binary: `go build -o adserver ./cmd/server`

## Project Structure

```
.
├── api/                    # API documentation and OpenAPI spec
├── cmd/                    # Application entrypoints
│   └── server/             # Main server application
├── internal/               # Private application code
│   ├── config/             # Configuration handling
│   ├── handler/            # HTTP handlers
│   ├── model/              # Data models
│   └── service/            # Business logic
├── docker-compose.yml      # Docker Compose configuration
├── Dockerfile              # Docker build configuration
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
└── README.md               # Project documentation
```

Good luck!