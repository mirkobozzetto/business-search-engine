import { NextResponse } from "next/server";
import { proxyGet } from "@/lib/api-proxy";

export async function GET(_request: Request, { params }: { params: Promise<{ identifier: string }> }) {
  const { identifier } = await params;
  const res = await proxyGet(`/companies/lookup/${identifier}`);
  const data = await res.json();
  return NextResponse.json(data, { status: res.status });
}
