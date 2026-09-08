import type { components } from "../../../shared/api/openapi";
import { formatIDR, formatPaymentDate } from "../../../shared/lib/format";
import { StatusBadge } from "./StatusBadge";

type Payment = components["schemas"]["Payment"];

export function PaymentTable({ payments }: { payments: Payment[] }) {
  return (
    <div className="overflow-x-auto border-y border-slate-200" role="region" aria-label="Payments table" tabIndex={0}>
      <table className="min-w-[760px] w-full border-collapse text-left">
        <thead className="bg-slate-50 text-xs font-semibold uppercase tracking-wide text-slate-500">
          <tr>
            <th scope="col" className="px-6 py-3.5">Payment ID</th>
            <th scope="col" className="px-6 py-3.5">Merchant name</th>
            <th scope="col" className="px-6 py-3.5">Date</th>
            <th scope="col" className="px-6 py-3.5 text-right">Amount</th>
            <th scope="col" className="px-6 py-3.5">Status</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-200 bg-white text-sm text-slate-700">
          {payments.map((payment) => (
            <tr key={payment.id} className="transition hover:bg-slate-50/70">
              <td className="whitespace-nowrap px-6 py-4 font-semibold text-slate-900">{payment.id}</td>
              <td className="px-6 py-4">{payment.merchant}</td>
              <td className="whitespace-nowrap px-6 py-4 text-slate-600">{formatPaymentDate(payment.created_at)}</td>
              <td className="whitespace-nowrap px-6 py-4 text-right font-semibold tabular-nums text-slate-900">
                {formatIDR(payment.amount)}
              </td>
              <td className="px-6 py-4"><StatusBadge status={payment.status} /></td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
