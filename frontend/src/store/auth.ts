import { create } from "zustand";
import type { SessionUser } from "@/lib/types";

interface AuthState {
  accessToken: string | null;
  user: SessionUser | null;
  isAuthenticated: boolean;
  setAuth: (accessToken: string, user: SessionUser) => void;
  setAccessToken: (token: string) => void;
  clearAuth: () => void;
}

export const useAuthStore = create<AuthState>()((set) => ({
  accessToken: null,
  user: null,
  isAuthenticated: false,

  setAuth(accessToken, user) {
    set({ accessToken, user, isAuthenticated: true });
  },

  setAccessToken(token) {
    set((state) => ({
      accessToken: token,
      isAuthenticated: state.user !== null,
    }));
  },

  clearAuth() {
    set({ accessToken: null, user: null, isAuthenticated: false });
  },
}));
