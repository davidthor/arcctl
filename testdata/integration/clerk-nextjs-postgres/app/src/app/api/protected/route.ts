import { auth } from "@clerk/nextjs/server";
import { Pool } from "pg";
import { NextResponse } from "next/server";

// Create a connection pool (reused across requests)
const pool = new Pool({
  connectionString: process.env.DATABASE_URL,
});

export async function GET() {
  try {
    // Check authentication via Clerk
    const { userId } = await auth();

    if (!userId) {
      return NextResponse.json(
        { error: "Unauthorized", message: "Authentication required" },
        { status: 401 }
      );
    }

    // Query the database to verify connection
    const result = await pool.query("SELECT NOW() as server_time");
    const serverTime = result.rows[0].server_time;

    return NextResponse.json({
      success: true,
      userId,
      dbConnected: true,
      serverTime,
      message: "Protected route accessed successfully",
    });
  } catch (error) {
    console.error("Error in protected route:", error);

    // Check if it's a database connection error
    if (error instanceof Error && error.message.includes("connect")) {
      return NextResponse.json(
        {
          error: "Database Error",
          message: "Failed to connect to database",
          dbConnected: false,
        },
        { status: 500 }
      );
    }

    return NextResponse.json(
      { error: "Internal Server Error", message: "An unexpected error occurred" },
      { status: 500 }
    );
  }
}
