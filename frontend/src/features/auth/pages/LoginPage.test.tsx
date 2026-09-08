import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { ApiError } from "../../../shared/api/error";
import { useAuthStore } from "../store/authStore";
import { LoginPage } from "./LoginPage";

vi.mock("../api/authApi", () => ({ login: vi.fn() }));

import { login } from "../api/authApi";

const mockedLogin = vi.mocked(login);

function renderLogin() {
  return render(
    <MemoryRouter initialEntries={["/login"]} future={{ v7_startTransition: true, v7_relativeSplatPath: true }}>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/dashboard" element={<p>Dashboard destination</p>} />
      </Routes>
    </MemoryRouter>
  );
}

describe("LoginPage", () => {
  beforeEach(() => {
    localStorage.clear();
    useAuthStore.setState({ token: null, email: null, role: null });
  });

  it("validates required fields without calling the API", async () => {
    renderLogin();
    await userEvent.click(screen.getByRole("button", { name: "Sign in" }));
    expect(screen.getByText("Enter your email and password.")).toBeInTheDocument();
    expect(mockedLogin).not.toHaveBeenCalled();
  });

  it("keeps the email and shows a safe credential error", async () => {
    mockedLogin.mockRejectedValue(new ApiError(401, { code: 401, message: "unauthorized" }));
    renderLogin();
    await userEvent.type(screen.getByLabelText("Email"), "cs@test.com");
    await userEvent.type(screen.getByLabelText("Password"), "wrong");
    await userEvent.click(screen.getByRole("button", { name: "Sign in" }));

    expect(await screen.findByText("The email or password is incorrect.")).toBeInTheDocument();
    expect(screen.getByLabelText("Email")).toHaveValue("cs@test.com");
  });

  it("disables only the submit action while login is pending", async () => {
    let resolveLogin!: (value: { token: string; email: string; role: "cs" }) => void;
    mockedLogin.mockReturnValue(new Promise((resolve) => {
      resolveLogin = resolve;
    }));
    renderLogin();
    await userEvent.type(screen.getByLabelText("Email"), "cs@test.com");
    await userEvent.type(screen.getByLabelText("Password"), "password");
    await userEvent.click(screen.getByRole("button", { name: "Sign in" }));

    expect(screen.getByRole("button", { name: "Signing in…" })).toBeDisabled();
    expect(screen.getByLabelText("Email")).toBeEnabled();
    expect(screen.getByLabelText("Password")).toBeEnabled();

    resolveLogin({ token: "token", email: "cs@test.com", role: "cs" });
    expect(await screen.findByText("Dashboard destination")).toBeInTheDocument();
  });

  it("stores the session and navigates after successful login", async () => {
    mockedLogin.mockResolvedValue({ token: "token-1", email: "cs@test.com", role: "cs" });
    renderLogin();
    await userEvent.type(screen.getByLabelText("Email"), " cs@test.com ");
    await userEvent.type(screen.getByLabelText("Password"), "password");
    await userEvent.click(screen.getByRole("button", { name: "Sign in" }));

    expect(await screen.findByText("Dashboard destination")).toBeInTheDocument();
    expect(mockedLogin).toHaveBeenCalledWith("cs@test.com", "password");
    await waitFor(() => expect(useAuthStore.getState().token).toBe("token-1"));
  });
});
