import api from "@/lib/api";
import { User } from "@/types/auth"

export async function register(username: string, email: string, password: string): Promise<User> {
    const res = await api.post("/auth/register", { username, email, password });
    return res.data.user;
}

export async function login(email: string, password: string): Promise<User> {
  const res = await api.post("/auth/login", { email, password });
  return res.data.user;
}

export async function logout(): Promise<void> {
  await api.post("/auth/logout");
  window.location.href = "/login";
}

export async function getMe(): Promise<User> {
  const res = await api.get("/auth/me");
  return res.data;
}
