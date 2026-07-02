import { Pool } from 'pg';

const pool = new Pool();

export class Invoice {
  async processAndPrint(id: string) {
    const result = await pool.query('SELECT * FROM invoices WHERE id = $1', [id]);
    const invoice = result.rows[0];
    let total = 0;
    for (let i = 0; i < invoice.items.length; i++) {
      if (invoice.items[i].type === 'standard') {
        total += invoice.items[i].price * 1.08;
      } else if (invoice.items[i].type === 'discounted') {
        total += invoice.items[i].price * 0.9 * 1.08;
      } else {
        total += invoice.items[i].price;
      }
    }
    console.log('Invoice ' + id + ' total: ' + total);
    return total;
  }
}
