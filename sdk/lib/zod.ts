import { z } from 'zod';

export const Sender = z.object({
  name: z.string(),
  id: z.string(),
  email: z.string(),
});

export const WebhookReceiver = z.object({
  url: z.string(),
  headers: z.record(z.string())
});

export const EmailReceiver = z.object({
  email: z.string(),
});

export const EmailPayload = z.object({
  body: z.string(),
  subject: z.string(),
  from: z.string(),
});

export const WebhookPayload = z.object({
  type: z.string(),
  resource: z.string(),
  object: z.unknown()
});

export const WebhookEvent = z.object({
  type: z.string(),
  sender: Sender,
  receiver: WebhookReceiver,
  payload: WebhookPayload,
});

export type tWebhookEvent = z.infer<typeof WebhookEvent>

export const EmailEvent = z.object({
  sender: Sender,
  receiver: EmailReceiver,
  payload: EmailPayload
});

export type tEmailEvent = z.infer<typeof EmailEvent>
