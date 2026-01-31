import { Pool } from "pg";
import { NextResponse } from "next/server";

// Health check endpoint for liveness probe
// This endpoint is NOT protected by Clerk auth

const pool = new Pool({
  connectionString: process.env.DATABASE_URL,
});

export async function GET() {
  try {
    // Check database connectivity
    const result = await pool.query("SELECT 1 as health_check");
    const dbHealthy = result.rows[0].health_check === 1;

    return NextResponse.json({
      status: "healthy",
      database: dbHealthy ? "connected" : "disconnected",
      timestamp: new Date().toISOString(),
    });
  } catch (error) {
    console.error("Health check failed:", error);

    return NextResponse.json(
      {
        status: "unhealthy",
        database: "disconnected",
        error: error instanceof Error ? error.message : "Unknown error",
        timestamp: new Date().toISOString(),
      },
      { status: 503 }
    );
  }
}
