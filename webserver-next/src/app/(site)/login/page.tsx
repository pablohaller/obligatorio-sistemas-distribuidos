"use client";

import { useState, useEffect, useCallback } from "react";
import { signIn, useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import Input from "@/app/components/ui/Input/Input";
import Button from "@/app/components/ui/Button/Button";
import { toast } from "react-toastify";
import {
  IconCircleCheck,
  IconLayoutDashboard,
  IconLoader3,
} from "@tabler/icons-react";
import { DASHBOARD_PATH, REGISTER_PATH } from "@/app/constants/routes";
import { FormData, defaultFormData } from "./helpers";
import {
  checkPayloadEmptyFiels,
  getFormDataErrors,
  getFormDataPayload,
  handleFormData,
} from "@/app/utils/helpers";
import { twMerge } from "tailwind-merge";
import LoadingSpinner from "@/app/components/LoadingSpinner/LoadingSpinner";

export default function Page() {
  const session = useSession();
  const router = useRouter();
  const [formData, setFormData] = useState<FormData>(defaultFormData);
  const [loading, setLoading] = useState<boolean>(false);
  const [successMessage, setSuccessMessage] = useState<boolean>(false);

  const handleFrom = useCallback(
    (key: string, attribute: string, value: unknown) => {
      setFormData(handleFormData(formData, key, attribute, value));
    },
    [formData]
  );

  useEffect(() => {
    if (session?.status === "authenticated") {
      if (successMessage) {
        setTimeout(() => router.push(DASHBOARD_PATH), 1000);
      } else {
        router.push(DASHBOARD_PATH);
      }
    }
  }, [session, successMessage, router]);

  const handleLogin = async () => {
    setLoading(true);
    const payload = getFormDataPayload(formData);
    const emptyFields = checkPayloadEmptyFiels(payload);
    if (emptyFields.length) {
      setFormData(getFormDataErrors(formData, emptyFields));
      setLoading(false);
      return;
    }
    const request = await signIn("credentials", {
      ...payload,
      redirect: false,
    });

    if (request?.error) {
      toast.error(request?.error, { theme: "colored" });
      setLoading(false);
      return;
    }

    setLoading(false);
    setSuccessMessage(true);
  };

  const {
    email: {
      label: emailLabel,
      value: emailValue,
      error: emailError,
      required: emailRequired,
    },
    password: {
      label: passwordLabel,
      value: passwordValue,
      error: passwordError,
      required: passwordRequired,
    },
  } = formData;

  return (
    <div className="h-screen grid place-items-center">
      <div className="   bg-white p-4 rounded-xl drop-shadow-md md:w-2/5 relative">
        {loading && <LoadingSpinner />}
        {successMessage && (
          <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-10 flex flex-col items-center text-center">
            <IconCircleCheck className="h-20 w-20 text-sky-500" />
            <div className="text-3xl">Usuario encontrado</div>
            <div>Redireccionando...</div>
          </div>
        )}
        <div
          className={twMerge(
            (loading || successMessage) && "blur-sm pointer-events-none"
          )}
        >
          <IconLayoutDashboard className="h-8 w-8 md:h-12 md:w-12 text-sky-500 m-auto" />
          <div className="font-rubik font-500 text-sky-500 text-center p-5 text-xl md:text-3xl">
            NEXTJS APP
          </div>
          <Input
            label={emailLabel}
            placeholder={emailLabel}
            required={emailRequired}
            value={emailValue}
            error={emailError}
            onChange={(e) => handleFrom("email", "value", e.target.value)}
          />
          <Input
            label={passwordLabel}
            type="password"
            placeholder={passwordLabel}
            required={passwordRequired}
            value={passwordValue}
            error={passwordError}
            onChange={(e) => handleFrom("password", "value", e.target.value)}
          />
          <div className="mt-4">
            <Button variant="contained" fullWidth onClick={handleLogin}>
              Iniciar sesión
            </Button>
            <Button
              variant="default"
              fullWidth
              onClick={() => router.push(REGISTER_PATH)}
            >
              Registrarse
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
