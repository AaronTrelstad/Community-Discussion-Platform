"use client"
import Link from "next/link"
import { useAuth } from "@/lib/context/AuthContext";

export default function Navbar() {
    const { user, logout, loading } = useAuth();

    return (
        <nav className="bg-white border-b border-gray-200 sticky top-0 z-50">
            <div className="max-w-5xl mx-auto px-4 h-12 flex items-center justify-between">
                <Link href="/" className="font-bold text-orange-500 text-lg">
                    Discuss
                </Link>

                <div className="flex items-center gap-3">
                    {loading ? (
                        <div className="w-20 h-8 bg-gray-100 rounded animate-pulse" />
                    ) : user ? (
                        <>
                            <span className="text-sm text-gray-600">u/{user.username}</span>
                            <button
                                onClick={logout}
                                className="text-sm text-gray-500 hover:text-gray-800"
                            >
                                Log out
                            </button>
                        </>
                    ) : (
                        <>
                            <Link
                                href="/login"
                                className="text-sm text-gray-600 hover:text-gray-900"
                            >
                                Log in
                            </Link>
                            <Link
                                href="/register"
                                className="text-sm bg-orange-500 text-white px-4 py-1.5 rounded-full hover:bg-orange-600 transition-colors"
                            >
                                Register
                            </Link>
                        </>
                    )}
                </div>
            </div>
        </nav>
    );
}
