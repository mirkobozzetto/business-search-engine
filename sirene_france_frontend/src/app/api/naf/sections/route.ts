import { NextResponse } from "next/server";
import { proxyGet } from "@/lib/api-proxy";

export async function GET() {
  const res = await proxyGet("/naf/sections");
  const data = await res.json();
  return NextResponse.json(data, { status: res.status });
}
