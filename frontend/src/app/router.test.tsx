import { render, screen } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { beforeEach, describe, expect, it } from "vitest";
import { useAuthStore } from "../features/auth/store/authStore";
import { ProtectedRoute } from "./router";

function renderProtectedRoute() {
  return render(
    <MemoryRouter
      initialEntries={["/dashboard"]}
      future={{ v7_startTransition: true, v7_relativeSplatPath: true }}
    >
      <Routes>
        <Route element={<ProtectedRoute />}>
          <Route path="/dashboard" element={<p>Dashboard</p>} />
        </Route>
        <Route path="/login" element={<p>Login</p>} />
      </Routes>
    </MemoryRouter>
  );
}

describe("ProtectedRoute", () => {
  beforeEach(() => {
    localStorage.clear();
    useAuthStore.setState({ token: null, email: null, role: null });
  });

  it("redirects an unauthenticated visitor to login", () => {
    renderProtectedRoute();
    expect(screen.getByText("Login")).toBeInTheDocument();
  });

  it("renders the protected route with a session", () => {
    useAuthStore.getState().setSession({ token: "token", email: "cs@test.com", role: "cs" });
    renderProtectedRoute();
    expect(screen.getByText("Dashboard")).toBeInTheDocument();
  });
});
