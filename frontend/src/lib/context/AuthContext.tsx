"use client"
import { createContext, useContext, useEffect, useState } from "react";
import { getMe, logout as logoutFn } from "@/lib/auth";
import { User } from "@/types/auth";

type AuthContextType = {
    user: User | null;
    setUser: (user: User) => void;
    logout: () => void;
    loading: boolean;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        const init = async () => {
            if (
                window.location.pathname.includes("/login") ||
                window.location.pathname.includes("/register")
            ) {
                setLoading(false);
                return;
            }

            try {
                const user = await getMe();
                setUser(user);
            } catch {
                setUser(null);
            } finally {
                setLoading(false);
            }
        };

        init();
    }, []);

    return (
        <AuthContext.Provider value={{ user, setUser, logout: logoutFn, loading }}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const ctx = useContext(AuthContext);
    if (!ctx) throw new Error("useAuth must be used within AuthProvider");
    return ctx;
}
