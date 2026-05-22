import { NextRequest, NextResponse } from "next/server";

const PUBLIC_PREFIXES = ["/api/auth", "/_next", "/favicon.ico"];

export function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl;

  if (PUBLIC_PREFIXES.some((p) => pathname.startsWith(p))) {
    return NextResponse.next();
  }

  const hasRefreshToken = req.cookies.has("refresh_token");
  const isAuthPage = pathname === "/login" || pathname === "/register";

  if (isAuthPage && hasRefreshToken) {
    // Authenticated users should not stay on login/register pages.
    const url = req.nextUrl.clone();
    url.pathname = "/inbox";
    return NextResponse.redirect(url);
  }

  if (!isAuthPage && !hasRefreshToken) {
    // Protected page, no session — send to login.
    const url = req.nextUrl.clone();
    url.pathname = "/login";
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!_next/static|_next/image|favicon.ico).*)"],
};
