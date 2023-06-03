"use client";

import { useState, useEffect } from "react";
import { signIn, useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import Input from "@/app/components/ui/Input/Input";
import Button from "@/app/components/ui/Button/Button";
import { toast } from "react-toastify";

interface FormData {
  username: string;
  email: string;
  password: string;
}

export default function Page() {
  const session = useSession();
  const router = useRouter();
  const [formData, setFormData] = useState<FormData>({
    username: "",
    email: "",
    password: "",
  });

  useEffect(() => {
    if (session?.status === "authenticated") {
      router.push("/dashboard");
    }
  }, [session]);

  const loginUser = async (e: any) => {
    e.preventDefault();
    signIn("credentials", { ...formData, redirect: false }).then((callback) => {
      if (callback?.error) {
        toast.error(callback?.error, { theme: "colored" });
      }

      if (callback?.ok && !callback?.error) {
        alert("Logged in");
      }
    });
  };

  const { email, password, username } = formData;

  return (
    <div className="bg-gradient-to-b from-cyan-50 to-sky-200 h-screen grid place-items-center">
      <div className=" bg-white p-4 rounded-xl drop-shadow-md md:w-1/4">
        <div className="font-rubik font-500 text-sky-500 text-center p-5 text-3xl">
          NextJS APP
        </div>
        <Input
          label="Nombre de usuario"
          placeholder="Nombre de Usuario"
          required
          value={username}
          onChange={(e) =>
            setFormData({ ...formData, username: e.target.value })
          }
        />
        <Input
          label="E-mail"
          placeholder="E-mail"
          required
          value={email}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
        />
        <Input
          label="Contraseña"
          placeholder="Contraseña"
          type="password"
          required
          value={password}
          onChange={(e) =>
            setFormData({ ...formData, password: e.target.value })
          }
        />
        <Button variant="contained" fullWidth onClick={loginUser}>
          Login
        </Button>
        <Button
          variant="default"
          fullWidth
          onClick={() =>
            toast.warning("Aún no implementado", { theme: "colored" })
          }
        >
          Registrarse
        </Button>
      </div>
    </div>
  );
}
