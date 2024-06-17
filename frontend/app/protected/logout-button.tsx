"use client";

import { LogoutAction } from "@/actions/auth";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { toast } from "sonner";

const LogoutButton = () => {
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  const handleLogout = async () => {
    setIsLoading(true);
    const response = await LogoutAction();
    if (response.success) {
      toast.success(response.success);
      setIsLoading(false);
      router.push("/login");
    } else {
      toast.error(response.error);
      setIsLoading(false);
    }
  };

  return (
    <Button disabled={isLoading} onClick={handleLogout}>
      Log Out
    </Button>
  );
};

export default LogoutButton;
