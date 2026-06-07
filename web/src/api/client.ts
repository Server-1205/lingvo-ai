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
  severity?: string;
  category?: string;
  learning_tip?: string;
  rule_violated?: string;
}

export interface PremiumAnalysis {
  overall_grade: string;
  strengths: string[];
  areas_for_improvement: string[];
  suggested_topic: string;
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
  premium_analysis?: PremiumAnalysis;
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
  ease_factor: number;
  interval: number;
  last_reviewed_at?: string | null;
  next_review: string | null;
  created_at: string;
}

export interface VocabLookupRequest {
  word: string;
}

export interface VocabLookupResponse {
  translation_uz: string;
  translation_ru: string;
  examples: string[];
  level: string;
}

export interface AddVocabRequest {
  word: string;
}

export interface GrammarRequest {
  text: string;
}

export interface GrammarResponse {
  corrections: Correction[];
}

export interface ProgressResponse {
  messages_sent: number;
  words_learned: number;
  quizzes_taken: number;
  streak_days: number;
  level: string;
}

export interface DailyProgressEntry {
  date: string;
  messages_sent: number;
  words_learned: number;
  quizzes_taken: number;
}

export interface ProgressHistoryResponse {
  entries: DailyProgressEntry[];
}

export interface QuizRequest {
  topic?: string;
  count?: number;
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
  questions: LevelQuestion[];
  level?: string;
}

export interface LevelQuestion {
  question: string;
  options: string[];
  correct: number;
  level: string;
  explanation_uz: string;
  explanation_ru: string;
}

export interface LevelSaveRequest {
  level: string;
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

export interface VocabListResponse {
  words: VocabWord[];
  total: number;
  due_count: number;
}

export async function getVocab(): Promise<VocabListResponse> {
  return request<VocabListResponse>('/api/vocab');
}

export async function addVocab(req: AddVocabRequest): Promise<VocabLookupResponse> {
  return request<VocabLookupResponse>('/api/vocab', {
    method: 'POST',
    body: JSON.stringify(req),
  });
}

export async function deleteVocab(id: number): Promise<{ status: string }> {
  return request<{ status: string }>(`/api/vocab/${id}`, {
    method: 'DELETE',
  });
}

export async function lookupVocab(req: VocabLookupRequest): Promise<VocabLookupResponse> {
  return request<VocabLookupResponse>('/api/vocab/lookup', {
    method: 'POST',
    body: JSON.stringify(req),
  });
}

export async function checkGrammar(req: GrammarRequest): Promise<GrammarResponse> {
  return request<GrammarResponse>('/api/grammar', {
    method: 'POST',
    body: JSON.stringify(req),
  });
}

export async function getProgress(): Promise<ProgressResponse> {
  return request<ProgressResponse>('/api/progress');
}

export async function getProgressHistory(days: number = 7): Promise<ProgressHistoryResponse> {
  return request<ProgressHistoryResponse>(`/api/progress/history?days=${days}`);
}

export async function getQuiz(req: QuizRequest): Promise<QuizResponse> {
  return request<QuizResponse>('/api/quiz', {
    method: 'POST',
    body: JSON.stringify(req),
  });
}

export async function saveLevel(req: LevelSaveRequest): Promise<{ status: string; level: string }> {
  return request<{ status: string; level: string }>('/api/level/save', {
    method: 'POST',
    body: JSON.stringify(req),
  });
}

export interface SSEEvent {
  type: 'token' | 'corrections' | 'usage' | 'done' | 'premium_analysis';
  data?: unknown;
}

export async function chatStream(
  req: ChatRequest,
  onToken: (text: string) => void,
  onCorrections: (corrections: Correction[]) => void,
  onUsage: (usage: Usage) => void,
  onDone: () => void,
  onPremiumAnalysis?: (analysis: PremiumAnalysis) => void,
): Promise<void> {
  const initDataRaw = initData.raw();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (initDataRaw) {
    headers['X-Telegram-Init-Data'] = initDataRaw;
  }

  console.debug('[api] POST /api/chat/stream');

  const res = await fetch(`${BASE}/api/chat/stream`, {
    method: 'POST',
    headers,
    body: JSON.stringify(req),
  });

  if (!res.ok) {
    const errData = await res.json().catch(() => ({}));
    throw new ApiError(
      (errData as Record<string, unknown>)?.message as string ||
      (errData as Record<string, unknown>)?.error as string ||
      'Unknown error',
      res.status,
      errData as Record<string, unknown>,
    );
  }

  const reader = res.body?.getReader();
  if (!reader) {
    throw new ApiError('no reader', 0, {});
  }

  const decoder = new TextDecoder();
  let buffer = '';

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;

    buffer += decoder.decode(value, { stream: true });
    const lines = buffer.split('\n');
    buffer = lines.pop() || '';

    for (const line of lines) {
      if (!line.startsWith('data: ')) continue;

      const jsonStr = line.slice(6);
      let evt: SSEEvent;
      try {
        evt = JSON.parse(jsonStr) as SSEEvent;
      } catch {
        console.debug('[chat/stream] failed to parse:', jsonStr);
        continue;
      }

      console.debug('[chat/stream] event:', evt.type, evt.data);

      switch (evt.type) {
        case 'token':
          onToken(evt.data as string);
          break;
        case 'corrections':
          onCorrections(evt.data as Correction[]);
          break;
        case 'usage':
          onUsage(evt.data as Usage);
          break;
        case 'premium_analysis':
          if (onPremiumAnalysis) {
            onPremiumAnalysis(evt.data as PremiumAnalysis);
          }
          break;
        case 'done':
          onDone();
          break;
      }
    }
  }
}

export async function getLevelTestQuestions(): Promise<LevelResponse> {
  return request<LevelResponse>('/api/level', {
    method: 'POST',
    body: '{}',
  });
}

export const ReviewQuality = {
  Again: 0,
  Hard: 2,
  Good: 4,
  Easy: 5,
} as const;

export interface ReviewResponse {
  next_review: string;
}

export async function getDueWords(limit?: number): Promise<VocabWord[]> {
  const params = limit ? `?limit=${limit}` : '';
  return request<VocabWord[]>('/api/vocab/review' + params);
}

export async function submitReview(wordId: number, quality: number): Promise<ReviewResponse> {
  return request<ReviewResponse>('/api/vocab/review', {
    method: 'POST',
    body: JSON.stringify({ word_id: wordId, quality }),
  });
}

export { ApiError };
