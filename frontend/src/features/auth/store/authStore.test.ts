import { beforeEach, describe, expect, it } from "vitest";
import { useAuthStore } from "./authStore";

describe("auth store", () => {
  beforeEach(() => {
    localStorage.clear();
    useAuthStore.setState({ token: null, email: null, role: null });
  });

  it("stores and clears the minimal session", () => {
    useAuthStore.getState().setSession({
      token: "token-1",
      email: "cs@test.com",
      role: "cs"
    });

    expect(useAuthStore.getState()).toMatchObject({
      token: "token-1",
      email: "cs@test.com",
      role: "cs"
    });
    expect(localStorage.getItem("durianpay-auth")).toContain("token-1");

    useAuthStore.getState().clearSession();

    expect(useAuthStore.getState()).toMatchObject({
      token: null,
      email: null,
      role: null
    });
  });
});
