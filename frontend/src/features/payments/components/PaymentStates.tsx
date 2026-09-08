export function LoadingState() {
  return (
    <div className="border-y border-slate-200 px-6 py-16 text-center" role="status">
      <div className="mx-auto h-6 w-6 animate-spin rounded-full border-2 border-slate-200 border-t-primary" />
      <p className="mt-3 text-sm font-medium text-slate-600">Loading payments…</p>
    </div>
  );
}

export function ErrorState({ message, onRetry }: { message: string; onRetry: () => void }) {
  return (
    <div className="border-y border-slate-200 px-6 py-14 text-center" role="alert">
      <h2 className="font-semibold text-slate-900">Payments could not be loaded</h2>
      <p className="mx-auto mt-2 max-w-md text-sm text-slate-600">{message}</p>
      <button
        type="button"
        onClick={onRetry}
        className="mt-5 rounded-md bg-primary px-4 py-2 text-sm font-semibold text-white hover:bg-purple-800 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2"
      >
        Try again
      </button>
    </div>
  );
}

export function EmptyState({ filtered }: { filtered: boolean }) {
  return (
    <div className="border-y border-slate-200 px-6 py-14 text-center">
      <h2 className="font-semibold text-slate-900">{filtered ? "No matching payments" : "No payments yet"}</h2>
      <p className="mt-2 text-sm text-slate-600">
        {filtered ? "Choose another status to see more results." : "Payments will appear here when they are available."}
      </p>
    </div>
  );
}
