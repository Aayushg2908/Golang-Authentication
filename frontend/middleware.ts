import { cookies } from "next/headers";
import { DEFAULT_REDIRECT, authRoutes, publicRoutes } from "./constants";
import { NextRequest, NextResponse } from "next/server";

export async function middleware(request: NextRequest) {
  const token = cookies().get("token");
  const path = request.nextUrl.pathname;

  const isAuthRoute = authRoutes.includes(path);
  const isPublicRoute = publicRoutes.includes(path);

  if (isAuthRoute) {
    if (token) {
      return NextResponse.redirect(new URL(DEFAULT_REDIRECT, request.nextUrl));
    }

    return null;
  }

  if (!token && !isPublicRoute) {
    return NextResponse.redirect(new URL("/login", request.nextUrl));
  }

  return null;
}

export const config = {
  matcher: ["/((?!.+\\.[\\w]+$|_next).*)", "/", "/(api|trpc)(.*)"],
};
