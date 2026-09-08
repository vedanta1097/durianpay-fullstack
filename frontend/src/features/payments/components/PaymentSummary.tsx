type PaymentSummaryProps = {
  total: number;
  completed: number;
  processing: number;
  failed: number;
};

const summaryItems: Array<{ key: keyof PaymentSummaryProps; label: string }> = [
  { key: "total", label: "Total payments" },
  { key: "completed", label: "Completed" },
  { key: "processing", label: "Processing" },
  { key: "failed", label: "Failed" }
];

export function PaymentSummary(props: PaymentSummaryProps) {
  return (
    <section aria-label="Payment summary" className="grid grid-cols-2 border-y border-slate-200 lg:grid-cols-4">
      {summaryItems.map((item, index) => (
        <div
          key={item.key}
          className={`px-4 py-5 sm:px-6 ${index % 2 !== 0 ? "border-l border-slate-200" : ""} ${index > 1 ? "border-t border-slate-200 lg:border-t-0" : ""} ${index === 2 ? "lg:border-l" : ""}`}
        >
          <p className="text-sm text-slate-500">{item.label}</p>
          <p className="mt-1 text-2xl font-bold tabular-nums text-slate-900">{props[item.key]}</p>
        </div>
      ))}
    </section>
  );
}
