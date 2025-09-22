<div align="center">

# [![PongHub](imgs/band.png)](https://health.ch3nyang.top)

🌏 [Live Demo](https://health.ch3nyang.top) | 📖 [简体中文](README_CN.md)

</div>

## Introduction

PongHub is an open-source service status monitoring website designed to help users track and verify service availability. It supports:

- **🕵️ Zero-intrusion Monitoring** - Full-featured monitoring without code changes
- **🚀 One-click Deployment** - Automatically built with GitHub Actions, deployed to GitHub Pages
- **🌐 Cross-platform Support** - Compatible with public services like OpenAI and private deployments
- **🔍 Multi-port Detection** - Monitor multiple ports for a single service
- **🤖 Intelligent Response Validation** - Precise matching of status codes and regex validation of response bodies
- **🛠️ Custom Request Engine** - Flexible configuration of request headers/bodies, timeouts, and retry strategies
- **🔒 SSL Certificate Monitoring** - Automatic detection of SSL certificate expiration and notifications
- **📊 Real-time Status Display** - Intuitive service response time and status records
- **⚠️ Exception Alert Notifications** - Exception alert notifications using GitHub Actions

![Browser Screenshot](imgs/browser.png)

## Quick Start

1. Star and Fork [PongHub](https://github.com/WCY-dt/ponghub)

2. Modify the [`config.yaml`](config.yaml) file in the root directory to configure your service checks.

3. Modify the [`CNAME`](CNAME) file in the root directory to set your custom domain name.
   
   > If you do not need a custom domain, you can delete the `CNAME` file.

4. Commit and push your changes to your repository. GitHub Actions will automatically run and deploy to GitHub Pages and require no intervention.

> [!TIP]
> By default, GitHub Actions runs every 30 minutes. If you need to change the frequency, modify the `cron` expression in the [`.github/workflows/deploy.yml`](.github/workflows/deploy.yml) file.
> 
> Please do not set the frequency too high to avoid triggering GitHub's rate limits.

> [!IMPORTANT]
> If GitHub Actions does not trigger automatically, you can manually trigger it once.
> 
> Please ensure that GitHub Pages is enabled and that you have granted notification permissions for GitHub Actions.

## Configuration Guide

### Basic Configuration

The `config.yaml` file follows this format:

| Field                               | Type    | Description                                              | Required | Notes                                         |
|-------------------------------------|---------|----------------------------------------------------------|----------|-----------------------------------------------|
| `timeout`                           | Integer | Timeout for each request in seconds                      | ✖️       | Units are seconds, default is 5 seconds       |
| `max_retry_times`                   | Integer | Number of retries on request failure                     | ✖️       | Default is 2 retries                          |
| `max_log_days`                      | Integer | Number of days to retain logs                            | ✖️       | Default is 3 days                             |
| `cert_notify_days`                  | Integer | Days before SSL certificate expiration to notify         | ✖️       | Default is 7 days                             |
| `services`                          | Array   | List of services to monitor                              | ✔️       |                                               |
| `services.name`                     | String  | Name of the service                                      | ✔️       |                                               |
| `services.endpoints`                | Array   | List of endpoints to check for the service               | ✔️       |                                               |                                               |
| `services.endpoints.url`            | String  | URL to request                                           | ✔️       |                                               |
| `services.endpoints.method`         | String  | HTTP method for the request                              | ✖️       | Supports `GET`/`POST`/`PUT`, default is `GET` |
| `services.endpoints.headers`        | Object  | Request headers                                          | ✖️       | Key-value                                     |
| `services.endpoints.body`           | String  | Request body content                                     | ✖️       | Used only for `POST`/`PUT` requests           |
| `services.endpoints.status_code`    | Integer | Expected HTTP status code in response (default is `200`) | ✖️       | Default is `200`                              |
| `services.endpoints.response_regex` | String  | Regex to match the response body content                 | ✖️       |                                               |

Here is an example configuration file:

```yaml
timeout: 5
max_retry_times: 2
max_log_days: 3
cert_notify_days: 7
services:
  - name: "GitHub API"
    endpoints:
      - url: "https://api.github.com"
      - url: "https://api.github.com/repos/wcy-dt/ponghub"
        method: "GET"
        headers:
          Content-Type: application/json
          Authorization: Bearer your_token
        status_code: 200
        response_regex: "full_name"
  - name: "Example Website"
    endpoints:
      - url: "https://example.com/health"
        response_regex: "status"
      - url: "https://example.com/status"
        method: "POST"
        body: '{"key": "value"}'
```

### Special Parameters

ponghub now supports powerful parameterized configuration functionality, allowing the use of various types of dynamic variables in configuration files. These variables are generated and resolved in real-time during program execution.

#### 📅 Date and Time Parameters

Use the `{{%format}}` format to define date and time parameters:

- `{{%Y-%m-%d}}` - Current date, format: 2006-01-02 (e.g., 2025-09-22)
- `{{%H:%M:%S}}` - Current time, format: 15:04:05 (e.g., 17:30:45)
- `{{%s}}` - Unix timestamp (e.g., 1727859600)
- `{{%Y}}` - Current year (e.g., 2025)
- `{{%m}}` - Current month, format: 01-12
- `{{%d}}` - Current day, format: 01-31
- `{{%H}}` - Current hour, format: 00-23
- `{{%M}}` - Current minute, format: 00-59
- `{{%S}}` - Current second, format: 00-59
- `{{%B}}` - Full month name (e.g., September)
- `{{%b}}` - Short month name (e.g., Sep)
- `{{%A}}` - Full weekday name (e.g., Monday)
- `{{%a}}` - Short weekday name (e.g., Mon)

#### 🎲 Random Number Parameters

- `{{rand}}` - Generates a random number in the range 0–1000000
- `{{rand_int}}` - Generates a large-range random integer
- `{{rand(min,max)}}` - Generates a random number within a specified range
    - Example: `{{rand(1,100)}}` - Generates a random number between 1 and 100
    - Example: `{{rand(1000,9999)}}` - Generates a 4-digit random number

#### 🔤 Random String Parameters

- `{{rand_str}}` - Generates an 8-character random string (letters + numbers)
- `{{rand_str(length)}}` - Generates a random string of specified length
    - Example: `{{rand_str(16)}}` - Generates a 16-character random string
- `{{rand_str_secure}}` - Generates a 16-character cryptographically secure random string
- `{{rand_hex(length)}}` - Generates a random hexadecimal string of specified length
    - Example: `{{rand_hex(8)}}` - Generates an 8-character hexadecimal string
    - Example: `{{rand_hex(32)}}` - Generates a 32-character hexadecimal string

#### 🆔 UUID Parameters

- `{{uuid}}` - Generates a standard UUID (with hyphens)
    - Example: `bf3655f7-8a93-4822-a458-2913a6fe4722`
- `{{uuid_short}}- Generates a short UUID (without hyphens)
    - Example: `14d44b7334014484bb81b015fb2401bf`

#### 🌍 Environment Variable Parameters

- `{{env(variable_name)}}` - Reads the value of an environment variable
    - Example: `{{env(API_KEY)}}` - Reads the API_KEY environment variable
    - Example: `{{env(VERSION)}}` - Reads the VERSION environment variable
    - If the environment variable does not exist, returns an empty string

Ensure that the environment variable is set in your GitHub repository settings under "Settings" -> "Secrets and variables" -> "Actions".

#### 📊 Serial Number and Hash Parameters

- `{{seq}}` - Sequence number based on the current time (6-digit number)
- `{{seq_daily}}` - Daily sequence number (seconds since midnight)
- `{{hash_short}}` - Short hash value (6-digit hexadecimal)
- `{{hash_md5_like}}` - MD5-style long hash value (32-digit hexadecimal)

Below is an example configuration file:

```yaml
services:
  - name: "Parameterized Service"
    endpoints:
        - url: "https://api.example.com/data?date={{%Y-%m-%d}}&rand={{rand(1,100)}}"
        - url: "https://api.example.com/submit"
          method: "POST"
          headers:
            Content-Type: application/json
            X-Request-ID: "{{uuid}}"
          body: '{"session_id": "{{rand_str(16)}}", "timestamp": "{{%s}}"}'
```

## Development

This project uses Makefile for local development and testing. You can run the project locally with the following command:

```bash
make run
```

## Disclaimer

[PongHub](https://github.com/WCY-dt/ponghub) is intended for personal learning and research only. The developers are not responsible for its usage or outcomes. Do not use it for commercial purposes or illegal activities.
