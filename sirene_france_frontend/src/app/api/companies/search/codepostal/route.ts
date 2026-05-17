import { NextRequest, NextResponse } from "next/server";
import { proxyGet } from "@/lib/api-proxy";

export async function GET(request: NextRequest) {
  const params = request.nextUrl.searchParams;
  const code = params.get("code") || params.get("q");
  if (!code) {
    return NextResponse.json({ success: false, error: "q parameter required" }, { status: 400 });
  }
  const backendParams = new URLSearchParams();
  backendParams.set("q", code);
  if (params.get("limit")) backendParams.set("limit", params.get("limit")!);
  if (params.get("offset")) backendParams.set("offset", params.get("offset")!);
  const res = await proxyGet("/companies/search/codepostal", backendParams);
  const data = await res.json();
  return NextResponse.json(data, { status: res.status });
}
