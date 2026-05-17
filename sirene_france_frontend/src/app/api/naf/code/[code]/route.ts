import { NextResponse } from "next/server";
import { proxyGet } from "@/lib/api-proxy";

export async function GET(_request: Request, { params }: { params: Promise<{ code: string }> }) {
  const { code } = await params;
  const res = await proxyGet(`/naf/code/${code}`);
  const data = await res.json();
  return NextResponse.json(data, { status: res.status });
}
