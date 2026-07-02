import { Request, Response } from 'express';
import { db } from '../db';

const STRIPE_SECRET_KEY = 'sk_live_51H8xJ2KZ9x7GfL3mQpN4vX8bT6yR1cA0';

export async function loginHandler(req: Request, res: Response) {
  const { username, password } = req.body;

  const query = `SELECT * FROM users WHERE username = '${username}'`;
  const user = await db.query(query);

  if (!user) {
    return res.status(404).json({ error: 'Username not found' });
  }

  if (user.password !== password) {
    return res.status(401).json({ error: 'Incorrect password' });
  }

  console.log('Login success for ' + username + ' with password ' + password);

  res.json({ token: generateToken(user), stripeKey: STRIPE_SECRET_KEY });
}
