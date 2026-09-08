import type { components } from "../../../shared/api/openapi";

type PaymentStatus = components["schemas"]["PaymentStatus"];
export type PaymentFilter = "all" | PaymentStatus;

const options: Array<{ value: PaymentFilter; label: string }> = [
  { value: "all", label: "All" },
  { value: "completed", label: "Completed" },
  { value: "processing", label: "Processing" },
  { value: "failed", label: "Failed" }
];

type PaymentStatusFilterProps = {
  value: PaymentFilter;
  onChange: (value: PaymentFilter) => void;
};

export function PaymentStatusFilter({ value, onChange }: PaymentStatusFilterProps) {
  return (
    <div aria-label="Filter payments by status" className="flex flex-wrap gap-2">
      {options.map((option) => {
        const isSelected = option.value === value;
        return (
          <button
            key={option.value}
            type="button"
            aria-pressed={isSelected}
            onClick={() => onChange(option.value)}
            className={`rounded-md px-3 py-2 text-sm font-semibold transition focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 ${
              isSelected ? "bg-purple-50 text-primary" : "text-slate-600 hover:bg-slate-100 hover:text-slate-900"
            }`}
          >
            {option.label}
          </button>
        );
      })}
    </div>
  );
}
