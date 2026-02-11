import type { APIResponse } from "@/types/api";

export async function fetchAPI<T>(path: string, params?: Record<string, string>): Promise<APIResponse<T>> {
  const url = new URL(path, window.location.origin);
  if (params) {
    Object.entries(params).forEach(([key, value]) => {
      if (value) url.searchParams.set(key, value);
    });
  }
  const res = await fetch(url.toString());
  if (!res.ok) {
    throw new Error(`API error: ${res.status}`);
  }
  return res.json();
}
