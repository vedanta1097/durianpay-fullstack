import { Navigate, Outlet, createBrowserRouter } from "react-router-dom";
import { LoginPage } from "../features/auth/pages/LoginPage";
import { useAuthStore } from "../features/auth/store/authStore";
import { DashboardPage } from "../features/payments/pages/DashboardPage";

export function ProtectedRoute() {
  const isAuthenticated = useAuthStore((state) => state.token !== null);
  return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />;
}

export const router = createBrowserRouter([
  { path: "/login", element: <LoginPage /> },
  {
    element: <ProtectedRoute />,
    children: [
      { path: "/dashboard", element: <DashboardPage /> }
    ]
  },
  { path: "*", element: <Navigate to="/dashboard" replace /> }
]);
