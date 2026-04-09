# Twilio E2E (SMS) with Go

This repository contains a minimal end-to-end Twilio SMS demo in Go:

- `cmd/send`: sends an outbound SMS via Twilio REST API.
- `cmd/server`: runs a small HTTP server that handles Twilio webhooks:
  - `POST /sms/incoming` (inbound SMS via your Twilio number webhook) 
  - `POST /status` (delivery status callbacks sent by Twilio)

## Prerequisites

- Go installed (`go version` should work)
- A Twilio account with:
  - an SMS-capable `From` phone number
  - a `To` phone number you can receive SMS on

## Environment variables

Copy `.env.example` to your environment (or export them manually):

- `TWILIO_ACCOUNT_SID`
- `TWILIO_AUTH_TOKEN`
- `TWILIO_FROM`
- `TWILIO_TO`
- `TWILIO_BODY`
- `TWILIO_STATUS_CALLBACK_URL` (must be publicly reachable by Twilio)
- `TWILIO_INCOMING_WEBHOOK_URL` (used as a reference in docs; webhook must be configured in Twilio console)
- `PORT` (for the local server)

Example:
```bash
export TWILIO_ACCOUNT_SID="ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
export TWILIO_AUTH_TOKEN="xxxxxxxxxxxxxxxxxxxx"
export TWILIO_FROM="+15017250604"
export TWILIO_TO="+15558675309"
export TWILIO_BODY="Hello from Go!"
```

## Run locally (webhooks require a public URL)
                                     
Twilio must be able to reach your server. For local development, use `ngrok` (or similar):

1. Start the server:
```bash
make run-server PORT=8080
```
2. (Optional) Quick webhook sanity checks locally:
```bash
# Inbound SMS endpoint (returns TwiML)
curl -s -X POST "http://localhost:8080/sms/incoming" \
  -d "From=+11111111111" \
  -d "To=+22222222222" \
  -d "Body=Hi from curl"

# Delivery status endpoint (just logs + returns ok)
curl -s -X POST "http://localhost:8080/status" \
  -d "MessageSid=SMxxxxxxxx" \
  -d "MessageStatus=delivered" \
  -d "From=+11111111111" \
  -d "To=+22222222222"
```
3. In a second terminal, expose it:
```bash
ngrok http 8080
```
4. Set callback URLs from the `ngrok` output:
```bash
export TWILIO_STATUS_CALLBACK_URL="https://<YOUR_NGROK_HOST>/status"
export TWILIO_INCOMING_WEBHOOK_URL="https://<YOUR_NGROK_HOST>/sms/incoming"
```

## Configure Twilio for inbound SMS

In the Twilio Console, update your `From` phone number webhook (or messaging service webhook) to:

- URL: `TWILIO_INCOMING_WEBHOOK_URL`
- Method: `POST`

Now when you text your Twilio number, Twilio will call your server and you’ll get a TwiML response.

## Send an outbound SMS (and watch status callbacks)

Run:

```bash
make send
```

This sends a message and requests Twilio to POST delivery updates to `TWILIO_STATUS_CALLBACK_URL`.
Watch the server logs to see `queued`, `sent`, `delivered`, etc.

## Build

```bash
make build
```

