import { afterEach, describe, expect, it, vi } from "vitest";

describe("payments API", () => {
  afterEach(() => {
    vi.unstubAllGlobals();
    vi.resetModules();
  });

  it("uses mocked fetch and reads the latest bearer token", async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({
        payments: [],
        summary: { total: 30, completed: 13, processing: 9, failed: 8 }
      }), {
        status: 200,
        headers: { "Content-Type": "application/json" }
      })
    );
    vi.stubGlobal("fetch", fetchMock);
    const { useAuthStore } = await import("../../auth/store/authStore");
    const { listPayments } = await import("./paymentsApi");
    useAuthStore.setState({ token: "current-token", email: "cs@test.com", role: "cs" });

    await expect(listPayments("failed")).resolves.toEqual({
      payments: [],
      summary: { total: 30, completed: 13, processing: 9, failed: 8 }
    });

    const request = fetchMock.mock.calls[0][0] as Request;
    expect(request.headers.get("Authorization")).toBe("Bearer current-token");
    expect(request.url).toContain("status=failed");
  });
});
