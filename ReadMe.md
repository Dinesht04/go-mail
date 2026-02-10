# GoMail-Service

A high-performance, distributed **Email Microservice** written in Go. It leverages Redis for reliable task queuing and background processing, supporting everything from one-time OTPs to scheduled newsletter subscriptions.

**System Architecture**

![alt text](image.png)

## What's the need for this?

Sending emails is a "slow" I/O operation. If your main application waits for an SMTP server to respond during a user's request, the user experiences lag.

**GoMail-Service** solves this by:

1. **Decoupling:** Your main app drops a "task" into the queue and returns a response instantly.
2. **Scheduling:** Handles complex logic like recurring subscriptions (daily/weekly/monthly) using internal cron workers.
3. **Reliability:** Includes built-in retries and structured logging to ensure no email is lost.

## Tech Stack

* **Framework:** [Gin Gonic](https://github.com/gin-gonic/gin) (HTTP Routing)
* **Queue/Storage:** [go-redis](https://github.com/redis/go-redis) (Redis client for Go)
* **Scheduling:** [robfig/cron](https://github.com/robfig/cron) (For recurring tasks)
* **Email:** [Native Go SMTP](https://pkg.go.dev/net/smtp) (Built-in email support)
* **Logging:** [slog](https://pkg.go.dev/log/slog) (Structured JSON logging)

---

## API Endpoints & Task Types

### 1. The Task Producer (`/task` - POST)

The primary entry point. It accepts a task and enqueues it for background processing.

**General Task Structure:**

```json
{
    "task": "User Login Verification",
    "type": "generateOtp",
    "retries": 3,
    "payload": {
       ...
    }
}

```

#### Supported Task Types

| Task Type | Description | Payload Requirements |
| --- | --- | --- |
| **message** | Sends an immediate email with specified subject and content. | `userId` (email), `subject`, `content`. |
| **generateOtp** | Generates an OTP, stores it in Redis for verification, and sends it to the user. | `userId` (email), `length` (int, min = 4, max = 8). |
| **subscribe** | Registers a recurring cron job (hourly, daily, weekly, monthly) for newsletters. | `userId`(email), `frequency`, `content_type`, `subject`, `content`. |
| **unsubscribe** | Removes a user's record from a specific cron-based subscription. | `userId`, `content_type`. |

**Sample Response:**

```json
{
    "status": true,
    "msg": "Task enqueued successfully"
}

```

---

### 2. OTP Verification (`/verify` - POST)

Validates the OTP stored in the Redis HashMap against the user's email.

**Request Payload:**

```json
{
    "userEmail": "user@example.com",
    "otp": "123456"
}

```

**Response:**

```json
{
    "type": "otp verification",
    "verified": true | false
}

```

---

### 3. Subscription Management (`/subscriptionContent`)

Manage the content templates used for your recurring emails.

#### **POST (Create Content)**

Defines a new content type for subscriptions.

**Request Payload:**

```json
{
    "content_type": "Weekly_Newsletter",
    "subject": "Subject of a weekly newsletter",
    "content": "Content of a weekly newsletter",
}
```

**Response:**

```json
{
    "status": true, 
    "msg": "Content Type created successfully"
}
```

#### **PUT (Update Content)**

Modifies an existing template in the Redis HashMap.

```json
{
    "content_type": "Weekly_Newsletter",
    "subject": "New Subject of a weekly newsletter",
    "content": "New Content of a weekly newsletter",
}
```

**Response:**

```json
{
    "status": true, 
    "msg": "Content Type updated successfully"
}
```
---

### 4. Health & Metrics

#### **`/ping` (GET)**

* **Description:** Standard health check to verify server availability.
* **Response:** `{"message": "pong"}`

#### **`/metrics` (GET)**

* **Description:** Retrieves real-time execution statistics directly from Redis.
* **Response:**

```json
{
    "status": true,
    "Total Jobs Executed": 150,
    "Jobs Successful": 148,
    "Jobs Failed": 2
}

```

---

## Logging & Monitoring

The service uses `slog` to provide **structured JSON logging** to both the standard output stream and an `app.log` file. I plan to integrate it with ElasticSearch in the future.

## **Request Parameters Reference**

The `/task` endpoint accepts a JSON object with a top-level `Task` structure containing a nested `Payload`.

#### **Main Task Wrapper**

| Field | Type | Required | Description / Constraints |
| --- | --- | --- | --- |
| `task` | `string` | **Yes** | A descriptive name for the task. |
| `type` | `string` | **Yes** | Must be one of: `generateOtp`, `message`, `subscribe`, `unsubscribe`. |
| `retries` | `int` | **Yes** | The number of times the system should retry the task on failure. |
| `payload` | `object` | **Yes** | The nested object containing task-specific data. |

---

#### **Payload Object**

The requirements for these fields change dynamically based on the `type` selected in the main task wrapper.

| Field | Type | Required | Constraint / Possible Values |
| --- | --- | --- | --- |
| `userId` | `string` | **Yes** | Must be a valid email address. |
| `content_type` | `string` | Conditional | **Required** if type is `subscribe` or `unsubscribe`. |
| `length` | `int` | Conditional | **Required** if type is `generateOtp`. Must be between **4 and 8**. |
| `frequency` | `string` | Conditional | **Required** if type is `subscribe`. Values: `@hourly`, `@daily`, `@weekly`, `@monthly`. |
| `content` | `string` | Conditional | **Required** if type is `message`. The body of the email. |
| `subject` | `string` | Conditional | **Required** if type is `message`. The subject line of the email. |

---

Made by [Dinesh](https://x.com/Dinesht_04)
