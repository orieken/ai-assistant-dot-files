// src/pricing/discount-calculator.ts
// Deliberately has no tests — unit-tester fixture

export type UserTier = 'free' | 'pro' | 'enterprise';

export function calculateDiscount(
  basePrice: number,
  tier: UserTier,
  couponCode?: string
): number {
  let discount = 0;

  if (tier === 'pro') {
    discount = 0.1;
  } else if (tier === 'enterprise') {
    discount = 0.25;
  }

  if (couponCode === 'SAVE10') {
    discount += 0.1;
  } else if (couponCode === 'SAVE20') {
    discount += 0.2;
  }

  // Cap total discount at 40%
  if (discount > 0.4) {
    discount = 0.4;
  }

  return basePrice * (1 - discount);
}
