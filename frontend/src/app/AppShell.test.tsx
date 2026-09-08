import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { beforeEach, describe, expect, it } from "vitest";
import { useAuthStore } from "../features/auth/store/authStore";
import { AppShell } from "./AppShell";

describe("AppShell", () => {
  beforeEach(() => {
    localStorage.clear();
    useAuthStore.setState({ token: "token", email: "cs@test.com", role: "cs" });
  });

  it("clears the session and navigates to login when logging out", async () => {
    render(
      <MemoryRouter initialEntries={["/dashboard"]} future={{ v7_startTransition: true, v7_relativeSplatPath: true }}>
        <Routes>
          <Route element={<AppShell />}>
            <Route path="/dashboard" element={<p>Dashboard content</p>} />
          </Route>
          <Route path="/login" element={<p>Login destination</p>} />
        </Routes>
      </MemoryRouter>
    );

    await userEvent.click(screen.getByRole("button", { name: "Log out" }));

    expect(await screen.findByText("Login destination")).toBeInTheDocument();
    await waitFor(() => expect(useAuthStore.getState().token).toBeNull());
  });
});
