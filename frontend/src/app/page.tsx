"use client";
import { useAuth } from "@/lib/context/AuthContext";
import Link from "next/link";

export default function HomePage() {
    const { user, loading } = useAuth();

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="w-6 h-6 border-2 border-orange-500 border-t-transparent rounded-full animate-spin" />
            </div>
        );
    }

    return (
        <main className="max-w-5xl mx-auto px-4 py-8">
            {user ? (
                <div>
                    <h1 className="text-2xl font-bold mb-1">Welcome back, u/{user.username}</h1>
                    <p className="text-gray-500 text-sm">Your feed will appear here.</p>
                </div>
            ) : (
                <div className="text-center py-20">
                    <h1 className="text-4xl font-bold mb-4">Welcome to Discuss</h1>
                    <p className="text-gray-500 mb-8">Join communities, share ideas, and vote on content.</p>
                    <div className="flex gap-3 justify-center">
                        <Link
                            href="/register"
                            className="bg-orange-500 text-white px-6 py-2 rounded-full font-medium hover:bg-orange-600 transition-colors"
                        >
                            Get started
                        </Link>
                        <Link
                            href="/login"
                            className="border border-gray-300 px-6 py-2 rounded-full font-medium hover:bg-gray-50 transition-colors"
                        >
                            Log in
                        </Link>
                    </div>
                </div>
            )}
        </main>
    );
}
