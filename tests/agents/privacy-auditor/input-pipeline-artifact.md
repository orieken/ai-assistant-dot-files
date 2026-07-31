# Implementation Notes: Payment Webhook Handler

## Summary
Added Stripe webhook endpoint at `POST /webhooks/stripe`. Handler verifies the Stripe signature
and processes `payment_intent.succeeded` events.

## Stripe Test Credentials Used During Development
For local testing, the following test API key was used: `sk_test_4eC39HqLyjWDarjtT1zdp7dc`
This is a Stripe test key (not production) but was hardcoded in `src/webhooks/stripe.handler.ts`
during development and not yet removed.

## User Data Handled
The handler extracts `customer_email` and `customer_name` from the Stripe event payload and logs
them at `INFO` level for debugging:
```typescript
logger.info({ customer_email, customer_name, amount }, 'Payment succeeded');
```

## Environment Variables
`STRIPE_WEBHOOK_SECRET` is read from `process.env` — correct.
`STRIPE_API_KEY` is hardcoded in `src/config/stripe.config.ts` as `sk_test_4eC39HqLyjWDarjtT1zdp7dc`.
