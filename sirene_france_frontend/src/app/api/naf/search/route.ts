import { NextRequest, NextResponse } from "next/server";
import { proxyGet } from "@/lib/api-proxy";

export async function GET(request: NextRequest) {
  const params = request.nextUrl.searchParams;
  const res = await proxyGet("/naf/search", params);
  const data = await res.json();
  return NextResponse.json(data, { status: res.status });
}
