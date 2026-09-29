# Quote Project Backend

A RESTful backend service for managing businesses, customers, and quotes.

The backend is built with **Go**, **PostgreSQL**, and **Firebase Authentication**, with a strong focus on authentication, business-scoped authorization, clean separation of responsibilities, and secure resource ownership.

---

## Overview

Quote Project is designed as a business-oriented quotation management system.

Each authenticated user is associated with a business and can manage resources belonging to that business, including:

- Business profile
- Customers
- Quotes
- Quote items
- Quote pricing
- Quote statuses

The backend uses Firebase Authentication to identify users and PostgreSQL to store application data.

A core security principle of the project is:

> The backend never trusts a client-provided `firebase_uid` or `business_id`.

Instead, the authenticated Firebase identity is verified by the backend and used to resolve the corresponding business.

---

# Tech Stack

### Backend

- Go
- `net/http`
- REST API
- Firebase Admin SDK for Go

### Authentication

- Firebase Authentication
- Phone Authentication
- Firebase ID Tokens
- Google Application Default Credentials (ADC)

### Database

- PostgreSQL
- SQL
- Transaction-based quote operations
- Foreign key relationships

### Development / Infrastructure

- Docker for local PostgreSQL
- Google Cloud CLI
- Application Default Credentials
- Git
- GitHub
- Environment variables
- Vite-based Firebase Phone Auth development utility

---

# Architecture

The backend follows a layered structure:

```text
HTTP Request
     ↓
Authentication Middleware
     ↓
Business Context
     ↓
Handler
     ↓
Service
     ↓
Repository
     ↓
PostgreSQL
```

Each layer has a specific responsibility.

### Handler

Responsible for:

- HTTP request parsing
- URL/query parameter handling
- Reading authenticated context values
- JSON decoding/encoding
- HTTP status codes

### Service

Responsible for:

- Business logic
- Validation
- Pricing logic
- Status validation
- Resource relationship validation

### Repository

Responsible for:

- PostgreSQL queries
- Data persistence
- Transactions
- Resource ownership constraints

### Authentication Middleware

Responsible for:

- Reading the Authorization header
- Validating the Bearer token format
- Verifying Firebase ID Tokens
- Extracting the Firebase UID
- Storing the UID in the request context

### Business Context

Responsible for:

- Reading the authenticated Firebase UID
- Resolving the corresponding Business ID
- Storing the Business ID in the request context

---

# Authentication

Authentication is handled using **Firebase Authentication**.

The current authentication flow uses Firebase Phone Authentication.

The client authenticates with Firebase and receives a Firebase ID Token.

Protected backend requests include:

```http
Authorization: Bearer <FIREBASE_ID_TOKEN>
```

The Go backend verifies the token using the Firebase Admin SDK.

---

## Authentication Flow

```text
User
  ↓
Firebase Phone Authentication
  ↓
Firebase User
  ↓
Firebase ID Token
  ↓
Authorization: Bearer <ID Token>
  ↓
Go Authentication Middleware
  ↓
Firebase Admin VerifyIDToken()
  ↓
Firebase UID
  ↓
Request Context
```

An invalid, missing, or expired token results in an unauthorized request.

---

# Firebase Admin Authentication

Firebase Admin is initialized in the Go backend.

For local development, the project uses **Google Application Default Credentials (ADC)**.

This means the repository does **not** require a Firebase Service Account private-key JSON file.

This avoids committing long-lived private credentials into source control.

Local authentication can be configured using:

```bash
gcloud auth login
```

Set the Firebase / Google Cloud project:

```bash
gcloud config set project <PROJECT_ID>
```

Configure Application Default Credentials:

```bash
gcloud auth application-default login
```

The backend can then use those credentials through the Firebase Admin SDK.

---

# Business Context & Authorization

Authentication answers:

> Who is making this request?

Authorization must additionally answer:

> Which business is this authenticated user allowed to access?

After Firebase authentication succeeds, the Firebase UID is placed in the request context.

The Business Context then resolves:

```text
Firebase UID
     ↓
businesses.firebase_uid
     ↓
PostgreSQL
     ↓
Business ID
     ↓
Request Context
```

Protected business resources therefore follow this flow:

```text
Firebase ID Token
        ↓
Authenticate
        ↓
Firebase UID
        ↓
BusinessContext.Resolve
        ↓
Business ID
        ↓
Handler
        ↓
Service
        ↓
Repository
```

The current data model assumes:

```text
1 Firebase UID → 1 Business
```

This can be extended in the future if multi-business accounts are required.

---

# Authorization & Resource Ownership

Knowing a database resource ID is not enough to access that resource.

For example, a request for:

```text
GET /quotes/25
```

does not automatically authorize access to quote `25`.

The repository additionally verifies that the quote belongs to the authenticated business.

Conceptually:

```sql
SELECT ...
FROM quotes
WHERE id = $1
  AND business_id = $2;
```

The same principle is applied to customer operations.

This protects the API from cross-business resource access.

For example:

```text
Authenticated Business: 3
Requested Quote: 20
Quote Owner: Business 2
```

The effective database lookup becomes:

```text
quote_id = 20
AND
business_id = 3
```

No matching resource is returned.

The API therefore treats the resource as not found instead of exposing another business's data.

---

# Security Model

The backend follows several important security rules.

### Never trust client-provided Business IDs

The client does not decide which business it is operating on.

The Business ID comes from:

```text
Verified Firebase Token
→ Firebase UID
→ Business lookup
→ Business ID
```

### Never trust client-provided Firebase UIDs

When creating a business, the Firebase UID is obtained from the authenticated request context.

A UID supplied in the JSON request body is not trusted.

### Resource ownership is enforced in PostgreSQL queries

Customer and quote operations are scoped by both:

```text
Resource ID + Authenticated Business ID
```

### Authentication and authorization are separated

Authentication middleware handles Firebase identity.

Business Context handles application-level business identity.

Repositories enforce resource ownership.

---

# API

The API currently contains three primary modules:

```text
Businesses
Customers
Quotes
```

Protected routes require:

```http
Authorization: Bearer <FIREBASE_ID_TOKEN>
```

---

# Business API

## Create Business

```http
POST /businesses
```

Creates a business for the authenticated Firebase user.

The client does not need to provide `firebase_uid`.

Example request:

```json
{
  "business_name": "Example Business",
  "owner_name": "John Doe",
  "phone": "+972500000000",
  "email": "business@example.com",
  "address": "Example Address"
}
```

The backend automatically assigns:

```text
firebase_uid = authenticated Firebase UID
```

A Firebase user cannot create another business when one already exists under the current one-business-per-user model.

---

## Get Current Business

```http
GET /businesses/me
```

Returns the business associated with the authenticated Firebase user.

No Business ID is required from the client.

---

## Update Current Business

```http
PATCH /businesses/me
```

Updates the authenticated user's business.

Example:

```json
{
  "business_name": "Updated Business Name"
}
```

The Business ID is resolved from the authenticated request context.

---

# Customer API

All customer endpoints are scoped to the authenticated business.

## Create Customer

```http
POST /customers
```

Example:

```json
{
  "name": "John Customer",
  "phone": "+972501234567",
  "email": "customer@example.com",
  "address": "Customer Address",
  "notes": "Important customer"
}
```

The client does not provide the trusted Business ID.

The backend assigns it from the authenticated Business Context.

---

## Get Customers

```http
GET /customers
```

Returns customers belonging to the authenticated business.

---

## Get Customer

```http
GET /customers/{id}
```

Returns the customer only if it belongs to the authenticated business.

---

## Update Customer

```http
PATCH /customers/{id}
```

Example:

```json
{
  "name": "Updated Customer Name",
  "phone": "+972509999999"
}
```

Ownership is validated using both the customer ID and authenticated Business ID.

---

## Delete Customer

```http
DELETE /customers/{id}
```

Deletes the customer only within the authenticated business scope.

---

# Quote API

Quotes are also fully scoped to the authenticated business.

---

## Get Next Quote Number

```http
GET /quotes/next-number
```

Example response:

```json
{
  "quote_number": "Q-0001"
}
```

Quote numbering is business-specific.

This means:

```text
Business A → Q-0001
Business B → Q-0001
```

is valid.

---

## Create Quote

```http
POST /quotes
```

Example:

```json
{
  "customer_id": 4,
  "title": "Website Development",
  "description": "Website development quote",
  "pricing_method": "items",
  "additional_amount": "100",
  "discount_type": "percent",
  "discount_value": "10",
  "vat_rate": "18",
  "status": "draft",
  "notes": "First quote",
  "items": [
    {
      "description": "Website Design",
      "quantity": "2",
      "unit_price": "500",
      "total_overridden": false,
      "position": 1
    },
    {
      "description": "Development",
      "quantity": "3",
      "unit_price": "700",
      "total_overridden": false,
      "position": 2
    }
  ]
}
```

The Business ID is automatically assigned from the authenticated Business Context.

The backend also verifies that the selected customer belongs to the authenticated business.

---

## Get Quotes

```http
GET /quotes
```

Returns quotes belonging only to the authenticated business.

The endpoint supports period filtering.

Supported periods:

```text
today
week
month
all
```

Example:

```http
GET /quotes?period=month
```

---

## Get Quote

```http
GET /quotes/{id}
```

Returns a quote only when it belongs to the authenticated business.

Quote items are returned as part of the quote.

---

## Update Quote

```http
PUT /quotes/{id}
```

Updates an existing quote.

The current implementation treats quote updates as full quote updates.

Quote ownership is validated using:

```text
Quote ID + Authenticated Business ID
```

If the customer is changed, the backend verifies that the new customer also belongs to the authenticated business.

---

## Update Quote Status

```http
PATCH /quotes/{id}/status
```

Example:

```json
{
  "status": "sent"
}
```

Supported quote statuses:

```text
draft
sent
viewed
approved
rejected
expired
```

---

## Delete Quote

```http
DELETE /quotes/{id}
```

Deletes a quote only when it belongs to the authenticated business.

---

# Quote Pricing

The backend supports two pricing methods.

## Item-Based Pricing

```text
pricing_method = items
```

Each item contains:

```text
description
quantity
unit_price
total
position
```

The backend calculates item totals and quote totals.

Conceptually:

```text
Items Subtotal
      +
Additional Amount
      ↓
Subtotal
      -
Discount
      ↓
After Discount
      +
VAT
      ↓
Final Total
```

Example:

```text
Items Subtotal:     3100
Additional Amount:   100
                    ----
Subtotal:            3200

Discount 10%:        -320
                    ----
After Discount:      2880

VAT 18%:           +518.40
                    -------
Total:              3398.40
```

---

## Manual Pricing

```text
pricing_method = manual
```

Manual pricing allows a manually supplied subtotal instead of calculating the subtotal from quote items.

Pricing validation is handled by the service layer.

---

# Quote Items

Quote items currently belong to the parent quote and are managed through quote operations.

The current full quote update implementation replaces the stored quote items with the items supplied in the update request.

A future improvement is planned to support differential item operations:

```text
Create individual item
Update existing item
Delete individual item
```

instead of replacing all quote items during every full quote update.

---

# Error Handling

The API uses standard HTTP status codes depending on the operation.

Common examples include:

```text
200 OK
201 Created
204 No Content
400 Bad Request
401 Unauthorized
403 Forbidden
404 Not Found
409 Conflict
405 Method Not Allowed
500 Internal Server Error
```

Examples:

Missing authentication:

```text
401 Unauthorized
```

Invalid/expired Firebase token:

```text
401 Unauthorized
```

Business already exists:

```text
409 Conflict
```

Attempting to access another business's customer or quote:

```text
404 Not Found
```

Returning `404` for business-scoped resources avoids exposing whether another business owns a particular resource ID.

---

# Firebase Phone Authentication Development Tool

A small browser-based Firebase authentication utility is included under:

```text
tools/firebase-phone-test/
```

It is used for local development and backend authentication testing.

The tool uses Firebase's JavaScript SDK and Vite.

Install dependencies:

```bash
cd tools/firebase-phone-test
npm install
```

Start the development tool:

```bash
npm run dev
```

Firebase test phone credentials can be configured using environment variables.

Example structure:

```env
VITE_FIREBASE_API_KEY=...
VITE_FIREBASE_AUTH_DOMAIN=...
VITE_FIREBASE_PROJECT_ID=...
VITE_FIREBASE_APP_ID=...
VITE_TEST_PHONE_NUMBER=...
VITE_TEST_VERIFICATION_CODE=...
```

Do not commit this `.env` file.

---

# Environment & Git Security

Sensitive/local files are excluded from Git.

Example `.gitignore`:

```gitignore
.env
node_modules/
```

Do not commit:

- Firebase test credentials
- Environment files containing secrets
- Service Account private keys
- `node_modules`
- Database credentials

---

# Local Development

## Prerequisites

Install:

- Go
- PostgreSQL / Docker
- Google Cloud CLI
- Node.js and npm for the Firebase test utility
- Git

---

## Clone Repository

```bash
git clone <repository-url>
cd Quote-Project
```

---

## Install Go Dependencies

```bash
go mod download
```

---

## Configure Google Cloud Authentication

Login:

```bash
gcloud auth login
```

Configure the project:

```bash
gcloud config set project <PROJECT_ID>
```

Configure ADC:

```bash
gcloud auth application-default login
```

---

## Start PostgreSQL

The project uses PostgreSQL for persistent application data.

When using the project's Docker configuration, start the required services with the appropriate Docker Compose command for the repository configuration.

---

## Run the Backend

From the project root:

```bash
go run ./cmd/api
```

The local API runs on:

```text
http://localhost:8080
```

---

# Testing Protected Endpoints

First obtain a valid Firebase ID Token through Firebase Authentication.

Then send it with the request:

```bash
curl \
  -H "Authorization: Bearer <FIREBASE_ID_TOKEN>" \
  http://localhost:8080/businesses/me
```

For example:

```bash
curl \
  -H "Authorization: Bearer <FIREBASE_ID_TOKEN>" \
  http://localhost:8080/customers
```

The backend verifies the Firebase token and automatically resolves the Business ID.

---

# Authentication & Authorization Example

Consider this request:

```http
GET /quotes/6
Authorization: Bearer <ID_TOKEN>
```

Internally:

```text
1. Authentication middleware receives the request

2. Firebase verifies the ID Token

3. Firebase UID is extracted

4. UID is stored in request context

5. Business Context reads the UID

6. PostgreSQL resolves:
   Firebase UID → Business ID

7. Business ID is stored in request context

8. Quote handler receives:
   Quote ID = 6
   Business ID = authenticated business

9. Repository queries:
   WHERE id = 6
   AND business_id = authenticated_business_id

10. Quote is returned only if the business owns it
```

This prevents horizontal cross-business resource access.

---

# Project Structure

The backend is organized by domain and responsibility.

A simplified structure:

```text
Quote-Project/
│
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── auth/
│   │   ├── firebase.go
│   │   ├── middleware.go
│   │   └── business_context.go
│   │
│   ├── business/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── model.go
│   │   └── routes.go
│   │
│   ├── customer/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── model.go
│   │   └── routes.go
│   │
│   └── quote/
│       ├── handler.go
│       ├── service.go
│       ├── repository.go
│       ├── model.go
│       ├── routes.go
│       ├── pricing.go
│       └── status.go
│
├── tools/
│   └── firebase-phone-test/
│
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

The exact structure may evolve as additional modules are introduced.

---

# Current Request Flow

```text
                        ┌─────────────────────┐
                        │    Client / App     │
                        └──────────┬──────────┘
                                   │
                          Firebase Phone Auth
                                   │
                                   ▼
                        ┌─────────────────────┐
                        │ Firebase ID Token   │
                        └──────────┬──────────┘
                                   │
                         Authorization: Bearer
                                   │
                                   ▼
                        ┌─────────────────────┐
                        │ Auth Middleware     │
                        │ VerifyIDToken()     │
                        └──────────┬──────────┘
                                   │
                             Firebase UID
                                   │
                                   ▼
                        ┌─────────────────────┐
                        │ Business Context    │
                        └──────────┬──────────┘
                                   │
                              Business ID
                                   │
                                   ▼
                        ┌─────────────────────┐
                        │      Handler        │
                        └──────────┬──────────┘
                                   │
                                   ▼
                        ┌─────────────────────┐
                        │      Service        │
                        │ Business Logic      │
                        └──────────┬──────────┘
                                   │
                                   ▼
                        ┌─────────────────────┐
                        │     Repository      │
                        │ Ownership + SQL     │
                        └──────────┬──────────┘
                                   │
                                   ▼
                        ┌─────────────────────┐
                        │     PostgreSQL      │
                        └─────────────────────┘
```

---

# Current Features

- Firebase Phone Authentication
- Firebase ID Token verification
- Authentication middleware
- Firebase UID request context
- Firebase UID → Business ID resolution
- Business-scoped authorization
- Business profile management
- Customer management
- Quote management
- Quote items
- Automatic quote numbering per business
- Item-based pricing
- Manual pricing
- Discounts
- VAT calculation
- Quote status management
- Cross-business resource protection
- PostgreSQL persistence
- Docker-based local database development
- ADC-based Firebase Admin authentication
- Local Firebase Phone Auth test utility

---

# Future Improvements

Planned or possible future improvements include:

- React Native frontend integration
- Centralized frontend API client
- Automatic Firebase ID Token attachment
- Firebase token refresh handling
- Differential quote-item updates
- Automated unit and integration tests
- CI/CD pipeline
- Structured logging
- Centralized API error responses
- API documentation / OpenAPI
- Production deployment
- Cloud database deployment
- Multi-business user support if required
- Roles and permissions if multiple users are introduced per business

---

# Frontend Integration

The next development stage is a React Native frontend.

The expected integration flow will be:

```text
React Native
     ↓
Firebase Phone Authentication
     ↓
Firebase ID Token
     ↓
API Client
     ↓
Authorization: Bearer <ID Token>
     ↓
Go Backend
     ↓
Business Context
     ↓
Business / Customers / Quotes
```

The frontend should not manage or trust Business IDs for authorization.

Business identity remains a backend responsibility derived from the authenticated Firebase user.

---

# Development Status

The backend foundation is complete for the current project stage.

Implemented and manually verified:

- Authentication
- Business resolution
- Business authorization
- Business operations
- Customer operations
- Quote operations
- Quote pricing
- Quote status handling
- Cross-business access protection

The project is ready for frontend integration.

---

## License

This project is currently intended for development and educational use.