"use server";

import { LoginSchema, RegisterSchema } from "@/lib/schema";
import axios from "axios";
import { cookies } from "next/headers";
import { z } from "zod";

export const RegisterAction = async (
  values: z.infer<typeof RegisterSchema>
) => {
  try {
    const response = await axios.post("http://localhost:8080/api/register", {
      name: values.name,
      email: values.email,
      password: values.password,
    });
    if (response.status === 200) {
      return { success: "Account created successfully!" };
    }
  } catch (error: any) {
    if (axios.isAxiosError(error)) {
      return { error: error.response?.data.error };
    } else {
      return { error: "An unexpected error occurred. Please try again later." };
    }
  }
};

export const LoginAction = async (values: z.infer<typeof LoginSchema>) => {
  try {
    const response = await axios.post("http://localhost:8080/api/login", {
      email: values.email,
      password: values.password,
    });
    if (response.status === 200) {
      cookies().set("token", response.data.token, {
        maxAge: 60 * 60 * 24 * 7,
        sameSite: "none",
        secure: true,
      });
      return { success: "Logged in successfully!" };
    }
  } catch (error: any) {
    if (axios.isAxiosError(error)) {
      return { error: error.response?.data.error };
    } else {
      return { error: "An unexpected error occurred. Please try again later." };
    }
  }
};

export const getUser = async () => {
  try {
    const token = cookies().get("token");
    const response = await axios.get("http://localhost:8080/api/me", {
      headers: {
        Authorization: `Bearer ${token?.value}`,
      },
    });
    if (response.status === 200) {
      return { user: response.data.user };
    }
  } catch (error) {
    if (axios.isAxiosError(error)) {
      return { error: error.response?.data.error };
    } else {
      return { error: "An unexpected error occurred. Please try again later." };
    }
  }
};

export const LogoutAction = async () => {
  try {
    cookies().delete("token");
    return { success: "Logged out successfully!" };
  } catch (error: any) {
    return { error: "An unexpected error occurred. Please try again later." };
  }
};
