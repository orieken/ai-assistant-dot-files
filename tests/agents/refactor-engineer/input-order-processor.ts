// src/orders/order-processor.ts
// Deliberately violates complexity and Sandi Metz rules — refactor-engineer fixture.
// Cyclomatic complexity: processOrder = 9, applyPromotions = 8 (both above the < 7 threshold).
// Mixed concerns: validation, pricing, tax, inventory, and notification all in one class.
// No tests exist for this file.

import { Database } from '../infra/database';
import { EmailService } from '../infra/email-service';

export type OrderStatus = 'pending' | 'confirmed' | 'rejected' | 'fulfilled';
export type PaymentMethod = 'credit_card' | 'paypal' | 'bank_transfer';
export type CustomerTier = 'standard' | 'silver' | 'gold' | 'platinum';

export interface OrderItem {
  productId: string;
  quantity: number;
  unitPrice: number;
}

export interface Order {
  id: string;
  customerId: string;
  customerTier: CustomerTier;
  items: OrderItem[];
  paymentMethod: PaymentMethod;
  promoCode?: string;
  status: OrderStatus;
}

export class OrderProcessor {
  constructor(
    private db: Database,
    private email: EmailService
  ) {}

  // Cyclomatic complexity: 9 — too many nested conditionals
  processOrder(order: Order): { success: boolean; message: string; total: number } {
    if (!order.id || !order.customerId) {
      return { success: false, message: 'Missing order ID or customer ID', total: 0 };
    }
    if (order.items.length === 0) {
      return { success: false, message: 'Order has no items', total: 0 };
    }

    let subtotal = 0;
    for (const item of order.items) {
      if (item.quantity <= 0) {
        return { success: false, message: `Invalid quantity for product ${item.productId}`, total: 0 };
      }
      subtotal += item.quantity * item.unitPrice;
    }

    const discountedTotal = this.applyPromotions(subtotal, order.customerTier, order.promoCode);

    let taxRate = 0.08;
    if (order.paymentMethod === 'bank_transfer') {
      taxRate = 0.06;
    }
    const tax = discountedTotal * taxRate;
    const total = discountedTotal + tax;

    if (total > 10000 && order.paymentMethod === 'paypal') {
      return { success: false, message: 'PayPal not supported for orders over $10,000', total: 0 };
    }

    this.db.save('orders', { ...order, status: 'confirmed', total });
    this.email.send(order.customerId, `Order ${order.id} confirmed. Total: $${total.toFixed(2)}`);
    return { success: true, message: 'Order confirmed', total };
  }

  // Cyclomatic complexity: 8 — promo + tier logic tangled together
  applyPromotions(subtotal: number, tier: CustomerTier, promoCode?: string): number {
    let discount = 0;

    if (tier === 'silver') {
      discount = 0.05;
    } else if (tier === 'gold') {
      discount = 0.1;
    } else if (tier === 'platinum') {
      discount = 0.15;
    }

    if (promoCode === 'SUMMER10') {
      discount += 0.1;
    } else if (promoCode === 'FLASH20') {
      discount += 0.2;
    } else if (promoCode === 'VIP30') {
      if (tier === 'platinum') {
        discount += 0.3;
      } else {
        discount += 0.15;
      }
    }

    if (discount > 0.35) {
      discount = 0.35;
    }

    return subtotal * (1 - discount);
  }
}
