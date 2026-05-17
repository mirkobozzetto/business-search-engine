export interface Meta {
  count?: number;
  total?: number;
  limit?: number;
  offset?: number;
  page?: number;
  pages?: number;
  duration_ms?: number;
}

export interface APIResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  meta?: Meta;
}
