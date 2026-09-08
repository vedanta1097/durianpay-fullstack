import type { components } from "../../../shared/api/openapi";

type PaymentStatus = components["schemas"]["PaymentStatus"];

const styles: Record<PaymentStatus, string> = {
  completed: "bg-emerald-50 text-emerald-700 ring-emerald-600/20",
  processing: "bg-amber-50 text-amber-700 ring-amber-600/20",
  failed: "bg-red-50 text-red-700 ring-red-600/20"
};

export function StatusBadge({ status }: { status: PaymentStatus }) {
  return (
    <span className={`inline-flex rounded-md px-2 py-1 text-xs font-semibold capitalize ring-1 ring-inset ${styles[status]}`}>
      {status}
    </span>
  );
}
