import { clerkMiddleware, createRouteMatcher } from "@clerk/nextjs/server";

// Define protected routes - the /api/protected route requires authentication
const isProtectedRoute = createRouteMatcher(["/api/protected(.*)"]);

export default clerkMiddleware(async (auth, req) => {
  // Protect routes that match the protected pattern
  if (isProtectedRoute(req)) {
    await auth.protect();
  }
});

export const config = {
  // Match all routes except static files and _next
  matcher: [
    // Skip Next.js internals and all static files
    "/((?!_next|[^?]*\\.(?:html?|css|js(?!on)|jpe?g|webp|png|gif|svg|ttf|woff2?|ico|csv|docx?|xlsx?|zip|webmanifest)).*)",
    // Always run for API routes
    "/(api|trpc)(.*)",
  ],
};
