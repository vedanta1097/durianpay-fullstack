import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { ApiError } from "../../../shared/api/error";
import { useAuthStore } from "../../auth/store/authStore";
import { DashboardPage } from "./DashboardPage";

vi.mock("../api/paymentsApi", () => ({ listPayments: vi.fn() }));

import { listPayments } from "../api/paymentsApi";

const mockedListPayments = vi.mocked(listPayments);
const payments = [
  { id: "PAY-1", merchant: "Kopi Nusantara", status: "completed" as const, amount: 125000, created_at: "2026-03-10T08:15:00Z" },
  { id: "PAY-2", merchant: "Toko Sejahtera", status: "processing" as const, amount: 89000, created_at: "2026-03-10T07:30:00Z" },
  { id: "PAY-3", merchant: "Batik Indah", status: "failed" as const, amount: 210000, created_at: "2026-03-09T15:45:00Z" }
];
const summary = { total: 30, completed: 13, processing: 9, failed: 8 };
const allPaymentsResponse = { payments, summary };

function renderDashboard() {
  return render(
    <MemoryRouter initialEntries={["/dashboard"]} future={{ v7_startTransition: true, v7_relativeSplatPath: true }}>
      <Routes>
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/login" element={<p>Login destination</p>} />
      </Routes>
    </MemoryRouter>
  );
}

describe("DashboardPage", () => {
  beforeEach(() => {
    mockedListPayments.mockReset();
    useAuthStore.setState({ token: "token", email: "cs@test.com", role: "cs" });
  });

  it("renders the backend summary and requests filtered payments", async () => {
    mockedListPayments.mockResolvedValueOnce(allPaymentsResponse).mockResolvedValueOnce({ payments: [payments[2]], summary });
    renderDashboard();

    expect(await screen.findByText("PAY-1")).toBeInTheDocument();
    expect(screen.getByText("PAY-2")).toBeInTheDocument();
    expect(screen.getByText("PAY-3")).toBeInTheDocument();
    expect(screen.getByText("Total payments").nextElementSibling).toHaveTextContent("30");
    expect(screen.getByText("Completed", { selector: "p" }).nextElementSibling).toHaveTextContent("13");

    await userEvent.click(screen.getByRole("button", { name: "Failed" }));
    await waitFor(() => expect(mockedListPayments).toHaveBeenLastCalledWith("failed"));
    expect(await screen.findByText("PAY-3")).toBeInTheDocument();
    expect(screen.queryByText("PAY-1")).not.toBeInTheDocument();
    expect(screen.getByText("Total payments").nextElementSibling).toHaveTextContent("30");
  });

  it("shows an error and retries", async () => {
    mockedListPayments.mockRejectedValueOnce(new Error("offline")).mockResolvedValueOnce(allPaymentsResponse);
    renderDashboard();

    expect(await screen.findByText("Payments could not be loaded")).toBeInTheDocument();
    await userEvent.click(screen.getByRole("button", { name: "Try again" }));
    expect(await screen.findByText("PAY-1")).toBeInTheDocument();
    expect(mockedListPayments).toHaveBeenCalledTimes(2);
  });

  it("shows the initial empty state", async () => {
    mockedListPayments.mockResolvedValue({ payments: [], summary: { total: 0, completed: 0, processing: 0, failed: 0 } });
    renderDashboard();

    expect(await screen.findByText("No payments yet")).toBeInTheDocument();
  });

  it("shows the filtered empty state without changing the summary", async () => {
    const onePaymentSummary = { total: 1, completed: 1, processing: 0, failed: 0 };
    mockedListPayments
      .mockResolvedValueOnce({ payments: [payments[0]], summary: onePaymentSummary })
      .mockResolvedValueOnce({ payments: [], summary: onePaymentSummary });
    renderDashboard();

    expect(await screen.findByText("PAY-1")).toBeInTheDocument();
    await userEvent.click(screen.getByRole("button", { name: "Failed" }));

    expect(await screen.findByText("No matching payments")).toBeInTheDocument();
    expect(screen.getByText("Total payments").nextElementSibling).toHaveTextContent("1");
  });

  it("clears an expired session and navigates to login", async () => {
    mockedListPayments.mockRejectedValue(new ApiError(401, { code: 401, message: "unauthorized" }));
    renderDashboard();

    expect(await screen.findByText("Login destination")).toBeInTheDocument();
    await waitFor(() => expect(useAuthStore.getState().token).toBeNull());
  });
});
