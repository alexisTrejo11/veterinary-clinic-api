# API Schema

**API type:** REST

## Addresses

### `POST` /api/v2/me/addresses/

**Create my address (planned route)**

PLANNED—AddressHandler.CreateAddressForMe implemented; routes not registered.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | addresses |

#### Request body

**Content-Type:** `application/json`

**Schema (summary):**

```json
{
  "street": "string",
  "city": "string",
  "state": "string",
  "zip": "string",
  "country": "string",
  "is_default": "boolean"
}
```

**Example:**

```json
{
  "street": "123 Main St",
  "city": "Austin",
  "state": "TX",
  "zip": "78701",
  "country": "US",
  "is_default": true
}
```

#### Responses

- **201** — Address created

```json
{
  "success": true,
  "data": {
    "id": 7
  }
}
```

---

## Appointments

### `GET` /api/v2/me/appointments/

**List my appointments (customer)**

Customer sees only their own appointment requests and bookings.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | appointments |

#### Responses

- **200** — Appointment list

```json
{
  "success": true,
  "data": [
    {
      "id": 100,
      "pet_id": 5,
      "scheduled_date": "2026-06-15T14:00:00Z",
      "status": "pending"
    }
  ]
}
```

---

### `POST` /api/v2/me/appointments/

**Request appointment (customer)**

Customer requests a new visit; domain rules validate capacity and employee overlap when employee assigned.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | appointments |

#### Request body

**Content-Type:** `application/json`

**Schema (summary):**

```json
{
  "pet_id": "integer (required)",
  "scheduled_date": "datetime (required)",
  "reason": "string (optional)",
  "employee_id": "integer (optional)"
}
```

**Example:**

```json
{
  "pet_id": 5,
  "scheduled_date": "2026-06-15T14:00:00Z",
  "reason": "Annual checkup"
}
```

#### Responses

- **201** — Appointment requested

```json
{
  "success": true,
  "data": {
    "id": 100
  }
}
```

---

### `GET` /api/v2/employees/appointments/

**List my assigned appointments (employee)**

Authenticated employee or manager sees appointments assigned to them.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | appointments |

#### Responses

- **200** — Assigned appointments

```json
{
  "success": true,
  "data": []
}
```

---

### `GET` /api/v2/appointments/

**Search all appointments (admin)**

Clinic-wide appointment search and filters for admin/manager roles.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | appointments |

#### Responses

- **200** — Paginated appointments

```json
{
  "success": true,
  "data": []
}
```

---

### `POST` /api/v2/appointments/{id}/confirm

**Confirm appointment**

Manager/admin confirms pending appointment. Optional query employee_id.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | appointments |

#### Parameters

| Name | In | Type | Required | Description |
| --- | --- | --- | --- | --- |
| id | path | integer | Yes | Appointment ID |
| employee_id | query | integer | No | Assigning employee |

#### Responses

- **200** — Appointment confirmed

```json
{
  "success": true,
  "message": "Appointment has been updated successfully"
}
```

---

### `POST` /api/v2/appointments/{id}/cancel

**Cancel appointment**

Cancels appointment with optional reason query param.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | appointments |

#### Parameters

| Name | In | Type | Required | Description |
| --- | --- | --- | --- | --- |
| id | path | integer | Yes | Appointment ID |

#### Responses

- **200** — Appointment cancelled

```json
{
  "success": true,
  "message": "Appointment has been updated successfully"
}
```

---

## Auth

### `POST` /api/v2/auth/register

**Register a new user**

Creates admin, customer, or employee-linked user. Non-admin accounts may require email activation before login.

| | |
|---|---|
| **Auth required** | No |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | auth |

#### Request body

**Content-Type:** `application/json`

**Schema (summary):**

```json
{
  "email": "string (required)",
  "password": "string (required, min length from PASSWORD_MIN_LENGTH)",
  "role": "admin | customer | veterinarian | receptionist",
  "full_name": "string (required)"
}
```

**Example:**

```json
{
  "email": "owner@pets.example.com",
  "password": "SecurePass123!",
  "role": "customer",
  "full_name": "Maria Garcia"
}
```

#### Responses

- **201** — Registration accepted

```json
{
  "success": true,
  "message": "Registration successful",
  "data": "success registration. An Email will be sent to the user to activate their account.",
  "timestamp": "2026-06-01T12:00:00Z"
}
```

---

### `POST` /api/v2/auth/login

**Authenticate and obtain JWT**

Email/password login. Returns access and refresh tokens and user profile. May require 2FA verification when enabled on account.

| | |
|---|---|
| **Auth required** | No |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP; lockout after MAX_LOGIN_ATTEMPTS |
| **Tags** | auth |

#### Request body

**Content-Type:** `application/json`

**Schema (summary):**

```json
{
  "email": "string (required)",
  "password": "string (required)",
  "two_factor_code": "string (optional)"
}
```

**Example:**

```json
{
  "email": "vet@clinic.example.com",
  "password": "SecurePass123!"
}
```

#### Responses

- **200** — Login successful

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 2,
      "email": "vet@clinic.example.com",
      "role": "veterinarian",
      "full_name": "Dr. Ana Lopez"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

---

### `POST` /api/v2/auth/activate

**Activate account**

Activates pending user using user_id and activation code from email link (query or JSON body).

| | |
|---|---|
| **Auth required** | No |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | auth |

#### Parameters

| Name | In | Type | Required | Description |
| --- | --- | --- | --- | --- |
| user_id | query | string | Yes | User ID to activate |
| code | query | string | Yes | Activation code |

#### Responses

- **200** — Account activated

```json
{
  "success": true,
  "message": "Account activated successfully"
}
```

---

### `POST` /auth/refresh

**Refresh access token**

Issues new access token using valid refresh token. Requires Bearer authentication on refresh token flow per handler implementation.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | auth |

#### Request body

**Content-Type:** `application/json`

**Schema (summary):**

```json
{
  "refresh_token": "string (required)"
}
```

**Example:**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Responses

- **200** — New session issued

```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

---

### `POST` /auth/logout

**Logout current session**

Revokes refresh token for current device. Requires Authorization Bearer access token.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per user/IP |
| **Tags** | auth |

#### Responses

- **200** — Logged out

```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

---

### `POST` /auth/logout-all

**Logout all sessions**

Revokes all refresh tokens for the authenticated user.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per user/IP |
| **Tags** | auth |

#### Responses

- **200** — All sessions revoked

```json
{
  "success": true,
  "message": "All sessions logged out"
}
```

---

## Customers

### `GET` /api/v2/customers/

**Search customers**

Admin/manager paginated customer search with filters in request body (POST-style search via GET handler binding).

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | customers |

#### Responses

- **200** — Customer page

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "full_name": "Maria Garcia",
      "email": "owner@pets.example.com",
      "is_active": true
    }
  ]
}
```

---

### `POST` /api/v2/customers/

**Create customer**

Creates clinic customer record linked to user account when applicable.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | customers |

#### Request body

**Content-Type:** `application/json`

**Schema (summary):**

```json
{
  "user_id": "integer (optional)",
  "full_name": "string",
  "email": "string",
  "phone": "string"
}
```

**Example:**

```json
{
  "full_name": "John Petrov",
  "email": "john@example.com",
  "phone": "+15550001111"
}
```

#### Responses

- **201** — Customer created

```json
{
  "success": true,
  "data": {
    "id": 10
  }
}
```

---

### `GET` /api/v2/customers/{id}

**Get customer by ID**

Returns single customer for admin/manager.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | customers |

#### Parameters

| Name | In | Type | Required | Description |
| --- | --- | --- | --- | --- |
| id | path | integer | Yes | Customer ID |

#### Responses

- **200** — Customer found

```json
{
  "success": true,
  "data": {
    "id": 10,
    "full_name": "John Petrov"
  }
}
```

---

## Employees

### `GET` /api/v2/employees/

**Search employees**

Lists veterinary staff (veterinarians, receptionists) for scheduling and admin.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | employees |

#### Responses

- **200** — Employee page

```json
{
  "success": true,
  "data": [
    {
      "id": 3,
      "full_name": "Dr. Ana Lopez",
      "specialty": "general_practice"
    }
  ]
}
```

---

### `POST` /api/v2/employees/

**Create employee**

Registers new clinic employee profile.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | employees |

#### Request body

**Content-Type:** `application/json`

**Schema (summary):**

```json
{
  "full_name": "string",
  "specialty": "string (vet specialty enum)",
  "license_number": "string"
}
```

**Example:**

```json
{
  "full_name": "Dr. Sam Kim",
  "specialty": "surgery",
  "license_number": "VET-12345"
}
```

#### Responses

- **201** — Employee created

```json
{
  "success": true,
  "data": {
    "id": 4
  }
}
```

---

## Medical

### `GET` /api/v2/medical/sessions

**List clinical sessions (staff)**

Staff search medical sessions. Requires medical handler wired in bootstrap (currently nil—endpoint inactive until wired).

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | medical |

#### Responses

- **200** — Session list

```json
{
  "success": true,
  "data": [
    {
      "id": 50,
      "pet_id": 5,
      "diagnosis": "Routine wellness exam",
      "status": "completed"
    }
  ]
}
```

---

### `GET` /api/v2/me/medical/sessions

**List my medical sessions (customer)**

Customer read-only view of clinical sessions for their pets.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | medical |

#### Responses

- **200** — Customer sessions

```json
{
  "success": true,
  "data": []
}
```

---

## Notifications

### `GET` /api/v2/me/notifications

**List my notifications**

Authenticated users read their in-app notification inbox.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | notifications |

#### Responses

- **200** — Notification list

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "type": "appointment_reminder",
      "channel": "email",
      "read": false
    }
  ]
}
```

---

### `POST` /api/v2/notifications

**Send notification (staff)**

Employee/manager/admin manually triggers a notification to a user.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | notifications |

#### Request body

**Content-Type:** `application/json`

**Schema (summary):**

```json
{
  "user_id": "integer",
  "type": "string",
  "channel": "email | sms | in_app",
  "message": "string"
}
```

**Example:**

```json
{
  "user_id": 2,
  "type": "custom",
  "channel": "in_app",
  "message": "Your lab results are ready."
}
```

#### Responses

- **201** — Notification queued

```json
{
  "success": true,
  "data": {
    "id": 99
  }
}
```

---

## Payments

### `GET` /api/v2/payments/{id}

**Get payment by ID (planned route)**

PLANNED—PaymentHandler exists; mount under admin/manager group in router.go.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | payments |

#### Parameters

| Name | In | Type | Required | Description |
| --- | --- | --- | --- | --- |
| id | path | integer | Yes | Payment ID |

#### Responses

- **200** — Payment found

```json
{
  "success": true,
  "data": {
    "id": 20,
    "amount": 85.5,
    "status": "completed",
    "currency": "USD"
  }
}
```

---

## Pets

### `GET` /api/v2/pets/

**Search pets (planned route)**

PLANNED—PetHandler.SearchPets implemented but route registration pending in router.go. Admin/manager search.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | pets |

#### Responses

- **200** — Pet page (when mounted)

```json
{
  "success": true,
  "data": [
    {
      "id": 5,
      "name": "Luna",
      "species": "dog",
      "customer_id": 1
    }
  ]
}
```

---

### `GET` /api/v2/me/pets/

**List my pets (planned route)**

PLANNED—customer-scoped pet list via GetMyPets handler.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | pets |

#### Responses

- **200** — My pets

```json
{
  "success": true,
  "data": [
    {
      "id": 5,
      "name": "Luna",
      "species": "dog"
    }
  ]
}
```

---

## Profile

### `GET` /api/v2/profile/

**Get current user profile**

Returns profile for the authenticated user (any role).

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | profile |

#### Responses

- **200** — Profile retrieved

```json
{
  "success": true,
  "data": {
    "id": 2,
    "email": "owner@pets.example.com",
    "role": "customer",
    "full_name": "Maria Garcia",
    "phone": "+15551234567"
  }
}
```

---

### `PUT` /api/v2/profile/

**Update current user profile**

Updates editable fields on the authenticated user's profile.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | profile |

#### Request body

**Content-Type:** `application/json`

**Schema (summary):**

```json
{
  "full_name": "string (optional)",
  "phone": "string (optional)"
}
```

**Example:**

```json
{
  "full_name": "Maria G. Garcia",
  "phone": "+15559876543"
}
```

#### Responses

- **200** — Profile updated

```json
{
  "success": true,
  "message": "User has been updated successfully"
}
```

---

## Service

### `GET` /health

**Service health check**

Public liveness probe for ALB, Docker, and uptime monitors. Returns service name and version.

| | |
|---|---|
| **Auth required** | No |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP when enabled |
| **Tags** | service |

#### Responses

- **200** — Service is up

```json
{
  "status": "ok",
  "timestamp": 1717200000,
  "service": "clinic-vet-api",
  "version": "2.0.0"
}
```

---

### `GET` /swagger/index.html

**Swagger UI (development)**

Interactive OpenAPI UI when ENABLE_SWAGGER=true and ENVIRONMENT is not production. Base path /api/v2 in swag annotations.

| | |
|---|---|
| **Auth required** | No |
| **Rate limit** | Unlimited in dev |
| **Tags** | service |

#### Responses

- **200** — HTML Swagger UI

```json
{
  "note": "Open in browser at http://localhost:8000/swagger/index.html"
}
```

---

## Users

### `GET` /users/

**Search users (admin)**

Paginated user search. Requires JWT roles admin or manager (see role alignment note in docs).

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | users |

#### Parameters

| Name | In | Type | Required | Description |
| --- | --- | --- | --- | --- |
| page | query | integer | No | Page number |
| page_size | query | integer | No | Items per page |

#### Responses

- **200** — Paginated user list

```json
{
  "success": true,
  "data": [],
  "meta": {
    "pagination": {
      "current_page": 1,
      "page_size": 20,
      "total": 0,
      "total_pages": 0
    }
  }
}
```

---

### `POST` /users/

**Create user (admin)**

Admin creates a new user account with role assignment.

| | |
|---|---|
| **Auth required** | Yes |
| **Rate limit** | GLOBAL — RATE_LIMIT_RPS per IP |
| **Tags** | users |

#### Request body

**Content-Type:** `application/json`

**Schema (summary):**

```json
{
  "email": "string",
  "password": "string",
  "role": "string",
  "full_name": "string"
}
```

**Example:**

```json
{
  "email": "reception@clinic.example.com",
  "password": "TempPass123!",
  "role": "receptionist",
  "full_name": "Carlos Ruiz"
}
```

#### Responses

- **201** — User created

```json
{
  "success": true,
  "data": {
    "id": 15
  }
}
```

---

## Additional notes

# API Schema

> **Base URL (production placeholder):** `https://api.vet-clinic.example.com`

> **Auth header:** `Authorization: Bearer <access_token>` on protected routes.

> **Response shape:** All JSON APIs use `APIResponse` with `success`, `message`, `data`, optional `error` and `meta` (pagination), plus `timestamp`.

> **Path prefix inconsistency (important):** Public auth lives under `/api/v2/auth/*`, but logout, refresh, 2FA, and reset-password use `/auth/*` **without** `/api/v2`. User admin routes use `/users/*` without prefix. Configure ALB path rules and frontend clients accordingly.

> **Role guard mismatch (dangerous):** Router checks `manager` and `employee` roles; JWT issues `admin`, `customer`, `veterinarian`, `receptionist`. Staff may get 403 until roles are unified—treat as pre-production bug.

> **Inactive modules:** Medical routes skip registration when `MedicalHandler` is nil. Pets, payments, and addresses are documented as **planned** paths—handlers exist in code but are not mounted.

> **Rate limiting:** Global Gin middleware uses in-memory storage per IP (`RATE_LIMIT_ENABLED`, `RATE_LIMIT_RPS`, `RATE_LIMIT_BURST`). Not suitable for multi-instance ECS without Redis-backed storage upgrade.

> **Deploy probe:** Use `GET /health` for ALB health checks—not `/api/v2/health`.

