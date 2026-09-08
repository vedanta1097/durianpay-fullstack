import { Outlet, useNavigate } from "react-router-dom";
import { useAuthStore } from "../features/auth/store/authStore";

export function AppShell() {
  const email = useAuthStore((state) => state.email);
  const role = useAuthStore((state) => state.role);
  const clearSession = useAuthStore((state) => state.clearSession);
  const navigate = useNavigate();

  function handleLogout() {
    clearSession();
    navigate("/login", { replace: true });
  }

  return (
    <div className="min-h-screen bg-slate-50">
      <header className="border-b border-slate-200 bg-white">
        <div className="mx-auto flex max-w-screen-2xl items-center justify-between px-4 py-4 sm:px-6 lg:px-8">
          <div className="flex items-center gap-3">
            <img src="/logo.svg" alt="DurianPay" className="h-7 w-auto" />
            <div>
              <p className="text-xs text-slate-500">Internal dashboard</p>
            </div>
          </div>

          <div className="flex items-center gap-4">
            <div className="hidden text-right sm:block">
              <p className="text-sm font-medium text-slate-800">{email}</p>
              <p className="text-xs text-slate-500">
                {role === "cs" ? "CS" : role === "operation" ? "Operation" : ""}
              </p>
            </div>
            <button
              type="button"
              onClick={handleLogout}
              className="rounded-md border border-slate-300 px-3 py-2 text-sm font-semibold text-slate-700 transition hover:bg-slate-50 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2"
            >
              Log out
            </button>
          </div>
        </div>
      </header>
      <Outlet />
    </div>
  );
}
