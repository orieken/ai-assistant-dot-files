# Implementation Notes: Payment Webhook Handler

## What was built
Added `POST /webhooks/stripe` endpoint that receives Stripe events, verifies the signature,
and updates order status in the database.

## Files changed
- `src/payments/webhook-handler.ts` (new) — receives and routes Stripe webhook events
- `src/orders/order-service.ts` (updated) — added `markAsPaid(orderId)` method
- `src/payments/stripe-client.ts` (new) — thin wrapper around `stripe` npm package

## New env vars
- `STRIPE_WEBHOOK_SECRET` — required for signature verification

## Notes
- Stripe signature is verified via `stripe.webhooks.constructEvent()` before any processing.
- Unrecognized event types are silently dropped (200 response to avoid Stripe retries).
- No retry logic yet — if `markAsPaid` throws, the webhook returns 500 and Stripe will retry.

## Testing
- Unit tests added for signature verification and event routing.
- No integration test yet for the full webhook → order update flow.
