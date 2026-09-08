import { FormEvent, useState } from "react";
import { Navigate, useNavigate } from "react-router-dom";
import { ApiError } from "../../../shared/api/error";
import { login } from "../api/authApi";
import { useAuthStore } from "../store/authStore";

export function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const token = useAuthStore((state) => state.token);
  const setSession = useAuthStore((state) => state.setSession);
  const navigate = useNavigate();

  if (token !== null) {
    return <Navigate to="/dashboard" replace />;
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const normalizedEmail = email.trim();
    if (normalizedEmail === "" || password === "") {
      setError("Enter your email and password.");
      return;
    }

    setError(null);
    setIsSubmitting(true);
    try {
      const session = await login(normalizedEmail, password);
      setSession(session);
      navigate("/dashboard", { replace: true });
    } catch (requestError) {
      if (requestError instanceof ApiError && requestError.status === 401) {
        setError("The email or password is incorrect.");
      } else {
        setError("We could not sign you in. Check your connection and try again.");
      }
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <main className="grid min-h-screen place-items-center bg-slate-50 px-4 py-10">
      <section className="w-full max-w-md border border-slate-200 bg-white px-6 py-8 shadow-sm sm:px-8">
        <div className="mb-8">
          <img src="/logo.svg" alt="DurianPay" className="mb-6 h-8 w-auto" />
          <p className="text-sm font-medium text-primary">Internal dashboard</p>
          <h1 className="mt-2 text-2xl font-bold tracking-tight text-slate-900">Sign in to your account</h1>
          <p className="mt-2 text-sm leading-6 text-slate-600">Monitor incoming payments and their latest status.</p>
        </div>

        <form onSubmit={handleSubmit} noValidate>
          <div>
            <label htmlFor="email" className="block text-sm font-semibold text-slate-700">
              Email
            </label>
            <input
              id="email"
              name="email"
              type="email"
              autoComplete="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              className="mt-2 block w-full rounded-md border border-slate-300 px-3 py-2.5 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-primary focus:ring-2 focus:ring-purple-100"
              placeholder="you@company.com"
            />
          </div>

          <div className="mt-5">
            <label htmlFor="password" className="block text-sm font-semibold text-slate-700">
              Password
            </label>
            <input
              id="password"
              name="password"
              type="password"
              autoComplete="current-password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              className="mt-2 block w-full rounded-md border border-slate-300 px-3 py-2.5 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-primary focus:ring-2 focus:ring-purple-100"
              placeholder="Enter your password"
            />
          </div>

          <div className="mt-4 min-h-6" aria-live="polite">
            {error !== null && <p className="text-sm text-red-700">{error}</p>}
          </div>

          <button
            type="submit"
            disabled={isSubmitting}
            className="mt-3 w-full rounded-md bg-primary px-4 py-2.5 text-sm font-semibold text-white transition hover:bg-purple-800 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {isSubmitting ? "Signing in…" : "Sign in"}
          </button>
        </form>
      </section>
    </main>
  );
}
