import { useCallback, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import type { components } from "../../../shared/api/openapi";
import { ApiError } from "../../../shared/api/error";
import { useAuthStore } from "../../auth/store/authStore";
import { listPayments } from "../api/paymentsApi";
import { EmptyState, ErrorState, LoadingState } from "../components/PaymentStates";
import { PaymentStatusFilter, type PaymentFilter } from "../components/PaymentStatusFilter";
import { PaymentSummary } from "../components/PaymentSummary";
import { PaymentTable } from "../components/PaymentTable";

type Payment = components["schemas"]["Payment"];
type PaymentSummaryData = components["schemas"]["PaymentSummary"];
type RequestState = "loading" | "success" | "error";

const emptySummary: PaymentSummaryData = {
  total: 0,
  completed: 0,
  processing: 0,
  failed: 0
};

export function DashboardPage() {
  const [payments, setPayments] = useState<Payment[]>([]);
  const [summary, setSummary] = useState<PaymentSummaryData>(emptySummary);
  const [requestState, setRequestState] = useState<RequestState>("loading");
  const [requestError, setRequestError] = useState("");
  const [filter, setFilter] = useState<PaymentFilter>("all");
  const clearSession = useAuthStore((state) => state.clearSession);
  const navigate = useNavigate();

  const loadPayments = useCallback(async () => {
    setRequestState("loading");
    setRequestError("");
    try {
      const response = await listPayments(filter === "all" ? undefined : filter);
      setPayments(response.payments);
      setSummary(response.summary);
      setRequestState("success");
    } catch (error) {
      if (error instanceof ApiError && error.status === 401) {
        clearSession();
        navigate("/login", { replace: true });
        return;
      }
      setRequestError("Check your connection, then try again.");
      setRequestState("error");
    }
  }, [clearSession, filter, navigate]);

  useEffect(() => {
    void loadPayments();
  }, [loadPayments]);

  return (
    <main className="mx-auto max-w-screen-2xl px-4 py-8 sm:px-6 lg:px-8 lg:py-10">
      <div className="mb-7">
        <p className="text-sm font-semibold text-primary">Operations</p>
        <h1 className="mt-1 text-2xl font-bold tracking-tight text-slate-900 sm:text-3xl">Incoming payments</h1>
        <p className="mt-2 text-sm text-slate-600">Monitor recent payments and review their processing status.</p>
      </div>

      <div className="overflow-hidden border border-slate-200 bg-white shadow-sm">
        <PaymentSummary {...summary} />

        <div className="flex flex-col gap-3 px-4 py-4 sm:flex-row sm:items-center sm:justify-between sm:px-6">
          <PaymentStatusFilter value={filter} onChange={setFilter} />
          {requestState === "success" && (
            <p className="text-sm tabular-nums text-slate-500">
              Showing <span className="font-semibold text-slate-700">{payments.length}</span> of {summary.total}
            </p>
          )}
        </div>

        {requestState === "loading" && <LoadingState />}
        {requestState === "error" && <ErrorState message={requestError} onRetry={() => void loadPayments()} />}
        {requestState === "success" && payments.length === 0 && <EmptyState filtered={filter !== "all"} />}
        {requestState === "success" && payments.length > 0 && <PaymentTable payments={payments} />}
      </div>
    </main>
  );
}
