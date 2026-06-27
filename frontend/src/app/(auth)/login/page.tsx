"use client"
import { useState } from "react";
import { useRouter } from "next/navigation";
import { login } from "@/lib/auth";
import { useAuth } from "@/lib/context/AuthContext";
import { UserLogin } from "@/types/auth";
import Link from "next/link";
import axios from "axios";

export default function LoginPage() {
    const router = useRouter();
    const { setUser } = useAuth();
    const [form, setForm] = useState<UserLogin>({ email: "", password: "" });
    const [error, setError] = useState<string>("");
    const [loading, setLoading] = useState<boolean>(false);

    async function handleSubmit(e: React.FormEvent) {
        e.preventDefault();
        setLoading(true);
        setError("");

        try {
            const user = await login(form.email, form.password);
            setUser(user);
            router.push("/");
        } catch (err: unknown) {
            if (axios.isAxiosError(err)) {
                setError(err.response?.data?.error || "Login failed");
            } else {
                setError("An unexpected error occurred");
            }
        } finally {
            setLoading(false);
        }
    }

    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50">
            <div className="bg-white p-8 rounded-lg shadow w-full max-w-md">
                <h1 className="text-2xl font-bold mb-2">Welcome back</h1>
                <p className="text-gray-500 text-sm mb-6">Log in to your account</p>
                {error && (
                    <div className="bg-red-50 border border-red-200 text-red-600 text-sm px-4 py-3 rounded mb-4">
                        {error}
                    </div>
                )}
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
                        <input
                            type="email"
                            value={form.email}
                            onChange={(e) => setForm(prev => ({ ...prev, email: e.target.value }))}
                            className="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500"
                            required
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Password</label>
                        <input
                            type="password"
                            value={form.password}
                            onChange={(e) => setForm(prev => ({ ...prev, password: e.target.value }))}
                            className="w-full border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500"
                            required
                        />
                    </div>
                    <button
                        type="submit"
                        disabled={loading}
                        className="w-full bg-orange-500 text-white py-2 rounded font-medium hover:bg-orange-600 disabled:opacity-50 transition-colors"
                    >
                        {loading ? "Logging in..." : "Log in"}
                    </button>
                </form>
                <p className="mt-4 text-sm text-center text-gray-500">
                    Don&apos;t have an account?{" "}
                    <Link href="/register" className="text-orange-500 hover:underline font-medium">
                        Register
                    </Link>
                </p>
            </div>
        </div>
    );
}
