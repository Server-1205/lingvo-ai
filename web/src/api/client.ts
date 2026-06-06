import { initData } from '@telegram-apps/sdk';

const BASE = import.meta.env.VITE_API_URL || '';

export interface ChatRequest {
  text: string;
}

export interface Correction {
  original: string;
  corrected: string;
  explanation_uz: string;
  explanation_ru: string;
  type: string;
}

export interface Usage {
  daily_used: number;
  daily_limit: number;
  is_premium: boolean;
}

export interface ChatResponse {
  reply: string;
  corrections: Correction[];
  usage: Usage;
}

export interface SubscriptionResponse {
  active: boolean;
  plan?: string;
  expires_at?: string;
}

export interface InvoiceRequest {
  plan: string;
}

export interface InvoiceResponse {
  invoice_link: string;
  stars: number;
}

export interface VocabWord {
  id: number;
  word: string;
  translation: string;
  example: string;
  level: string;
  review_count: number;
  next_review: string | null;
}

export interface VocabListResponse {
  words: VocabWord[];
  total: number;
  due_count: number;
}

export interface ProgressResponse {
  level: string;
  messages: number;
  words: number;
  streak: number;
}

export interface QuizQuestion {
  question: string;
  options: string[];
  correct: number;
  explanation_uz: string;
  explanation_ru: string;
}

export interface QuizResponse {
  questions: QuizQuestion[];
}

export interface LevelResponse {
  level: string;
  score: number;
}

class ApiError extends Error {
  status: number;
  data: Record<string, unknown>;

  constructor(message: string, status: number, data: Record<string, unknown>) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.data = data;
  }
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const initDataRaw = initData.raw();
  const headers: Record<string, string> = {
    ...((options?.body || options?.method === 'POST') ? { 'Content-Type': 'application/json' } : {}),
  };
  if (initDataRaw) {
    headers['X-Telegram-Init-Data'] = initDataRaw;
  }

  console.debug(`[api] ${options?.method || 'GET'} ${path}`);

  const res = await fetch(`${BASE}${path}`, {
    ...options,
    headers: { ...headers, ...options?.headers },
  });

  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    console.debug(`[api] ${options?.method || 'GET'} ${path} → ${res.status}`, data);
    throw new ApiError(
      (data as Record<string, unknown>)?.message as string || (data as Record<string, unknown>)?.error as string || 'Unknown error',
      res.status,
      data as Record<string, unknown>,
    );
  }

  const data = await res.json() as T;
  console.debug(`[api] ${options?.method || 'GET'} ${path} → ${res.status}`, data);
  return data;
}

export async function chatSend(req: ChatRequest): Promise<ChatResponse> {
  return request<ChatResponse>('/api/chat', {
    method: 'POST',
    body: JSON.stringify(req),
  });
}

export async function getSubscription(): Promise<SubscriptionResponse> {
  return request<SubscriptionResponse>('/api/subscription');
}

export async function createInvoice(req: InvoiceRequest): Promise<InvoiceResponse> {
  return request<InvoiceResponse>('/api/create-invoice', {
    method: 'POST',
    body: JSON.stringify(req),
  });
}

export async function getVocab(page = 1, perPage = 20, dueOnly = false): Promise<VocabListResponse> {
  const params = new URLSearchParams({
    page: String(page),
    per_page: String(perPage),
  });
  if (dueOnly) params.set('due_only', 'true');
  return request<VocabListResponse>(`/api/vocab?${params}`);
}

export async function addVocab(word: string): Promise<VocabWord> {
  return request<VocabWord>('/api/vocab', {
    method: 'POST',
    body: JSON.stringify({ word }),
  });
}

export async function getProgress(): Promise<ProgressResponse> {
  return request<ProgressResponse>('/api/progress');
}

export { ApiError };
