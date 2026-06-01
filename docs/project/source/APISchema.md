---
type: "REST"

httpEndpoints:
  - id: "health-check"
    method: "GET"
    urlPath: "/health"
    summary: "Service health check"
    description: "Public liveness probe for ALB, Docker, and uptime monitors. Returns service name and version."
    tags: ["service"]
    authenticated: false
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP when enabled"
    responses:
      - status: 200
        description: "Service is up"
        example:
          status: "ok"
          timestamp: 1717200000
          service: "clinic-vet-api"
          version: "2.0.0"

  - id: "auth-register"
    method: "POST"
    urlPath: "/api/v2/auth/register"
    summary: "Register a new user"
    description: "Creates admin, customer, or employee-linked user. Non-admin accounts may require email activation before login."
    tags: ["auth"]
    authenticated: false
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    requestBody:
      contentType: "application/json"
      schema:
        email: "string (required)"
        password: "string (required, min length from PASSWORD_MIN_LENGTH)"
        role: "admin | customer | veterinarian | receptionist"
        full_name: "string (required)"
      example:
        email: "owner@pets.example.com"
        password: "SecurePass123!"
        role: "customer"
        full_name: "Maria Garcia"
    responses:
      - status: 201
        description: "Registration accepted"
        example:
          success: true
          message: "Registration successful"
          data: "success registration. An Email will be sent to the user to activate their account."
          timestamp: "2026-06-01T12:00:00Z"

  - id: "auth-login"
    method: "POST"
    urlPath: "/api/v2/auth/login"
    summary: "Authenticate and obtain JWT"
    description: "Email/password login. Returns access and refresh tokens and user profile. May require 2FA verification when enabled on account."
    tags: ["auth"]
    authenticated: false
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP; lockout after MAX_LOGIN_ATTEMPTS"
    requestBody:
      contentType: "application/json"
      schema:
        email: "string (required)"
        password: "string (required)"
        two_factor_code: "string (optional)"
      example:
        email: "vet@clinic.example.com"
        password: "SecurePass123!"
    responses:
      - status: 200
        description: "Login successful"
        example:
          success: true
          message: "Login successful"
          data:
            user:
              id: 2
              email: "vet@clinic.example.com"
              role: "veterinarian"
              full_name: "Dr. Ana Lopez"
            access_token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
            refresh_token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

  - id: "auth-activate"
    method: "POST"
    urlPath: "/api/v2/auth/activate"
    summary: "Activate account"
    description: "Activates pending user using user_id and activation code from email link (query or JSON body)."
    tags: ["auth"]
    authenticated: false
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    parameters:
      - name: "user_id"
        in: "query"
        type: "string"
        required: true
        description: "User ID to activate"
        example: "42"
      - name: "code"
        in: "query"
        type: "string"
        required: true
        description: "Activation code"
        example: "A1B2C3D4"
    responses:
      - status: 200
        description: "Account activated"
        example:
          success: true
          message: "Account activated successfully"

  - id: "auth-refresh"
    method: "POST"
    urlPath: "/auth/refresh"
    summary: "Refresh access token"
    description: "Issues new access token using valid refresh token. Requires Bearer authentication on refresh token flow per handler implementation."
    tags: ["auth"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    requestBody:
      contentType: "application/json"
      schema:
        refresh_token: "string (required)"
      example:
        refresh_token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    responses:
      - status: 200
        description: "New session issued"
        example:
          success: true
          data:
            access_token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
            refresh_token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

  - id: "auth-logout"
    method: "POST"
    urlPath: "/auth/logout"
    summary: "Logout current session"
    description: "Revokes refresh token for current device. Requires Authorization Bearer access token."
    tags: ["auth"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per user/IP"
    responses:
      - status: 200
        description: "Logged out"
        example:
          success: true
          message: "Logged out successfully"

  - id: "auth-logout-all"
    method: "POST"
    urlPath: "/auth/logout-all"
    summary: "Logout all sessions"
    description: "Revokes all refresh tokens for the authenticated user."
    tags: ["auth"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per user/IP"
    responses:
      - status: 200
        description: "All sessions revoked"
        example:
          success: true
          message: "All sessions logged out"

  - id: "profile-get"
    method: "GET"
    urlPath: "/api/v2/profile/"
    summary: "Get current user profile"
    description: "Returns profile for the authenticated user (any role)."
    tags: ["profile"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    responses:
      - status: 200
        description: "Profile retrieved"
        example:
          success: true
          data:
            id: 2
            email: "owner@pets.example.com"
            role: "customer"
            full_name: "Maria Garcia"
            phone: "+15551234567"

  - id: "profile-update"
    method: "PUT"
    urlPath: "/api/v2/profile/"
    summary: "Update current user profile"
    description: "Updates editable fields on the authenticated user's profile."
    tags: ["profile"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    requestBody:
      contentType: "application/json"
      schema:
        full_name: "string (optional)"
        phone: "string (optional)"
      example:
        full_name: "Maria G. Garcia"
        phone: "+15559876543"
    responses:
      - status: 200
        description: "Profile updated"
        example:
          success: true
          message: "User has been updated successfully"

  - id: "users-search"
    method: "GET"
    urlPath: "/users/"
    summary: "Search users (admin)"
    description: "Paginated user search. Requires JWT roles admin or manager (see role alignment note in docs)."
    tags: ["users"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    parameters:
      - name: "page"
        in: "query"
        type: "integer"
        required: false
        description: "Page number"
        example: 1
      - name: "page_size"
        in: "query"
        type: "integer"
        required: false
        description: "Items per page"
        example: 20
    responses:
      - status: 200
        description: "Paginated user list"
        example:
          success: true
          data: []
          meta:
            pagination:
              current_page: 1
              page_size: 20
              total: 0
              total_pages: 0

  - id: "users-create"
    method: "POST"
    urlPath: "/users/"
    summary: "Create user (admin)"
    description: "Admin creates a new user account with role assignment."
    tags: ["users"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    requestBody:
      contentType: "application/json"
      schema:
        email: "string"
        password: "string"
        role: "string"
        full_name: "string"
      example:
        email: "reception@clinic.example.com"
        password: "TempPass123!"
        role: "receptionist"
        full_name: "Carlos Ruiz"
    responses:
      - status: 201
        description: "User created"
        example:
          success: true
          data:
            id: 15

  - id: "customers-list"
    method: "GET"
    urlPath: "/api/v2/customers/"
    summary: "Search customers"
    description: "Admin/manager paginated customer search with filters in request body (POST-style search via GET handler binding)."
    tags: ["customers"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    responses:
      - status: 200
        description: "Customer page"
        example:
          success: true
          data:
            - id: 1
              full_name: "Maria Garcia"
              email: "owner@pets.example.com"
              is_active: true

  - id: "customers-create"
    method: "POST"
    urlPath: "/api/v2/customers/"
    summary: "Create customer"
    description: "Creates clinic customer record linked to user account when applicable."
    tags: ["customers"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    requestBody:
      contentType: "application/json"
      schema:
        user_id: "integer (optional)"
        full_name: "string"
        email: "string"
        phone: "string"
      example:
        full_name: "John Petrov"
        email: "john@example.com"
        phone: "+15550001111"
    responses:
      - status: 201
        description: "Customer created"
        example:
          success: true
          data:
            id: 10

  - id: "customers-get"
    method: "GET"
    urlPath: "/api/v2/customers/{id}"
    summary: "Get customer by ID"
    description: "Returns single customer for admin/manager."
    tags: ["customers"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    parameters:
      - name: "id"
        in: "path"
        type: "integer"
        required: true
        description: "Customer ID"
        example: 10
    responses:
      - status: 200
        description: "Customer found"
        example:
          success: true
          data:
            id: 10
            full_name: "John Petrov"

  - id: "employees-list"
    method: "GET"
    urlPath: "/api/v2/employees/"
    summary: "Search employees"
    description: "Lists veterinary staff (veterinarians, receptionists) for scheduling and admin."
    tags: ["employees"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    responses:
      - status: 200
        description: "Employee page"
        example:
          success: true
          data:
            - id: 3
              full_name: "Dr. Ana Lopez"
              specialty: "general_practice"

  - id: "employees-create"
    method: "POST"
    urlPath: "/api/v2/employees/"
    summary: "Create employee"
    description: "Registers new clinic employee profile."
    tags: ["employees"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    requestBody:
      contentType: "application/json"
      schema:
        full_name: "string"
        specialty: "string (vet specialty enum)"
        license_number: "string"
      example:
        full_name: "Dr. Sam Kim"
        specialty: "surgery"
        license_number: "VET-12345"
    responses:
      - status: 201
        description: "Employee created"
        example:
          success: true
          data:
            id: 4

  - id: "me-appointments-list"
    method: "GET"
    urlPath: "/api/v2/me/appointments/"
    summary: "List my appointments (customer)"
    description: "Customer sees only their own appointment requests and bookings."
    tags: ["appointments"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    responses:
      - status: 200
        description: "Appointment list"
        example:
          success: true
          data:
            - id: 100
              pet_id: 5
              scheduled_date: "2026-06-15T14:00:00Z"
              status: "pending"

  - id: "me-appointments-create"
    method: "POST"
    urlPath: "/api/v2/me/appointments/"
    summary: "Request appointment (customer)"
    description: "Customer requests a new visit; domain rules validate capacity and employee overlap when employee assigned."
    tags: ["appointments"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    requestBody:
      contentType: "application/json"
      schema:
        pet_id: "integer (required)"
        scheduled_date: "datetime (required)"
        reason: "string (optional)"
        employee_id: "integer (optional)"
      example:
        pet_id: 5
        scheduled_date: "2026-06-15T14:00:00Z"
        reason: "Annual checkup"
    responses:
      - status: 201
        description: "Appointment requested"
        example:
          success: true
          data:
            id: 100

  - id: "employee-appointments-list"
    method: "GET"
    urlPath: "/api/v2/employees/appointments/"
    summary: "List my assigned appointments (employee)"
    description: "Authenticated employee or manager sees appointments assigned to them."
    tags: ["appointments"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    responses:
      - status: 200
        description: "Assigned appointments"
        example:
          success: true
          data: []

  - id: "appointments-search"
    method: "GET"
    urlPath: "/api/v2/appointments/"
    summary: "Search all appointments (admin)"
    description: "Clinic-wide appointment search and filters for admin/manager roles."
    tags: ["appointments"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    responses:
      - status: 200
        description: "Paginated appointments"
        example:
          success: true
          data: []

  - id: "appointments-confirm"
    method: "POST"
    urlPath: "/api/v2/appointments/{id}/confirm"
    summary: "Confirm appointment"
    description: "Manager/admin confirms pending appointment. Optional query employee_id."
    tags: ["appointments"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    parameters:
      - name: "id"
        in: "path"
        type: "integer"
        required: true
        description: "Appointment ID"
        example: 100
      - name: "employee_id"
        in: "query"
        type: "integer"
        required: false
        description: "Assigning employee"
        example: 3
    responses:
      - status: 200
        description: "Appointment confirmed"
        example:
          success: true
          message: "Appointment has been updated successfully"

  - id: "appointments-cancel"
    method: "POST"
    urlPath: "/api/v2/appointments/{id}/cancel"
    summary: "Cancel appointment"
    description: "Cancels appointment with optional reason query param."
    tags: ["appointments"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    parameters:
      - name: "id"
        in: "path"
        type: "integer"
        required: true
        description: "Appointment ID"
        example: 100
    responses:
      - status: 200
        description: "Appointment cancelled"
        example:
          success: true
          message: "Appointment has been updated successfully"

  - id: "me-notifications-list"
    method: "GET"
    urlPath: "/api/v2/me/notifications"
    summary: "List my notifications"
    description: "Authenticated users read their in-app notification inbox."
    tags: ["notifications"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    responses:
      - status: 200
        description: "Notification list"
        example:
          success: true
          data:
            - id: 1
              type: "appointment_reminder"
              channel: "email"
              read: false

  - id: "notifications-send"
    method: "POST"
    urlPath: "/api/v2/notifications"
    summary: "Send notification (staff)"
    description: "Employee/manager/admin manually triggers a notification to a user."
    tags: ["notifications"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    requestBody:
      contentType: "application/json"
      schema:
        user_id: "integer"
        type: "string"
        channel: "email | sms | in_app"
        message: "string"
      example:
        user_id: 2
        type: "custom"
        channel: "in_app"
        message: "Your lab results are ready."
    responses:
      - status: 201
        description: "Notification queued"
        example:
          success: true
          data:
            id: 99

  - id: "medical-sessions-list"
    method: "GET"
    urlPath: "/api/v2/medical/sessions"
    summary: "List clinical sessions (staff)"
    description: "Staff search medical sessions. Requires medical handler wired in bootstrap (currently nil—endpoint inactive until wired)."
    tags: ["medical"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    responses:
      - status: 200
        description: "Session list"
        example:
          success: true
          data:
            - id: 50
              pet_id: 5
              diagnosis: "Routine wellness exam"
              status: "completed"

  - id: "me-medical-sessions"
    method: "GET"
    urlPath: "/api/v2/me/medical/sessions"
    summary: "List my medical sessions (customer)"
    description: "Customer read-only view of clinical sessions for their pets."
    tags: ["medical"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    responses:
      - status: 200
        description: "Customer sessions"
        example:
          success: true
          data: []

  - id: "pets-list-planned"
    method: "GET"
    urlPath: "/api/v2/pets/"
    summary: "Search pets (planned route)"
    description: "PLANNED—PetHandler.SearchPets implemented but route registration pending in router.go. Admin/manager search."
    tags: ["pets"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    responses:
      - status: 200
        description: "Pet page (when mounted)"
        example:
          success: true
          data:
            - id: 5
              name: "Luna"
              species: "dog"
              customer_id: 1

  - id: "me-pets-list-planned"
    method: "GET"
    urlPath: "/api/v2/me/pets/"
    summary: "List my pets (planned route)"
    description: "PLANNED—customer-scoped pet list via GetMyPets handler."
    tags: ["pets"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    responses:
      - status: 200
        description: "My pets"
        example:
          success: true
          data:
            - id: 5
              name: "Luna"
              species: "dog"

  - id: "payments-get-planned"
    method: "GET"
    urlPath: "/api/v2/payments/{id}"
    summary: "Get payment by ID (planned route)"
    description: "PLANNED—PaymentHandler exists; mount under admin/manager group in router.go."
    tags: ["payments"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    parameters:
      - name: "id"
        in: "path"
        type: "integer"
        required: true
        description: "Payment ID"
        example: 20
    responses:
      - status: 200
        description: "Payment found"
        example:
          success: true
          data:
            id: 20
            amount: 85.50
            status: "completed"
            currency: "USD"

  - id: "addresses-me-planned"
    method: "POST"
    urlPath: "/api/v2/me/addresses/"
    summary: "Create my address (planned route)"
    description: "PLANNED—AddressHandler.CreateAddressForMe implemented; routes not registered."
    tags: ["addresses"]
    authenticated: true
    rateLimit: "GLOBAL — RATE_LIMIT_RPS per IP"
    requestBody:
      contentType: "application/json"
      schema:
        street: "string"
        city: "string"
        state: "string"
        zip: "string"
        country: "string"
        is_default: "boolean"
      example:
        street: "123 Main St"
        city: "Austin"
        state: "TX"
        zip: "78701"
        country: "US"
        is_default: true
    responses:
      - status: 201
        description: "Address created"
        example:
          success: true
          data:
            id: 7

  - id: "swagger-ui"
    method: "GET"
    urlPath: "/swagger/index.html"
    summary: "Swagger UI (development)"
    description: "Interactive OpenAPI UI when ENABLE_SWAGGER=true and ENVIRONMENT is not production. Base path /api/v2 in swag annotations."
    tags: ["service"]
    authenticated: false
    rateLimit: "Unlimited in dev"
    responses:
      - status: 200
        description: "HTML Swagger UI"
        example:
          note: "Open in browser at http://localhost:8000/swagger/index.html"
---

# API Schema

> **Base URL (production placeholder):** `https://api.vet-clinic.example.com`

> **Auth header:** `Authorization: Bearer <access_token>` on protected routes.

> **Response shape:** All JSON APIs use `APIResponse` with `success`, `message`, `data`, optional `error` and `meta` (pagination), plus `timestamp`.

> **Path prefix inconsistency (important):** Public auth lives under `/api/v2/auth/*`, but logout, refresh, 2FA, and reset-password use `/auth/*` **without** `/api/v2`. User admin routes use `/users/*` without prefix. Configure ALB path rules and frontend clients accordingly.

> **Role guard mismatch (dangerous):** Router checks `manager` and `employee` roles; JWT issues `admin`, `customer`, `veterinarian`, `receptionist`. Staff may get 403 until roles are unified—treat as pre-production bug.

> **Inactive modules:** Medical routes skip registration when `MedicalHandler` is nil. Pets, payments, and addresses are documented as **planned** paths—handlers exist in code but are not mounted.

> **Rate limiting:** Global Gin middleware uses in-memory storage per IP (`RATE_LIMIT_ENABLED`, `RATE_LIMIT_RPS`, `RATE_LIMIT_BURST`). Not suitable for multi-instance ECS without Redis-backed storage upgrade.

> **Deploy probe:** Use `GET /health` for ALB health checks—not `/api/v2/health`.
