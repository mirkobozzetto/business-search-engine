import { BACKEND_URL } from "./constants";

export async function proxyGet(path: string, params?: URLSearchParams): Promise<Response> {
  const url = new URL(`/api${path}`, BACKEND_URL);
  if (params) {
    params.forEach((value, key) => {
      if (value) url.searchParams.set(key, value);
    });
  }
  const res = await fetch(url.toString(), { cache: "no-store" });
  return res;
}
