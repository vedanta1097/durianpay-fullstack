import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { components } from "../../../shared/api/openapi";

type UserRole = components["schemas"]["UserRole"];

type AuthSession = {
  token: string;
  email: string;
  role: UserRole;
};

type AuthState = {
  token: string | null;
  email: string | null;
  role: UserRole | null;
  setSession: (session: AuthSession) => void;
  clearSession: () => void;
};

const emptySession = {
  token: null,
  email: null,
  role: null
};

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      ...emptySession,
      setSession: (session) => set(session),
      clearSession: () => set(emptySession)
    }),
    {
      name: "durianpay-auth",
      partialize: (state) => ({
        token: state.token,
        email: state.email,
        role: state.role
      })
    }
  )
);
